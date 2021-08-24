package m28to29

import (
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/badgerhelpers"
	"github.com/stackrox/rox/migrator/bolthelpers"
	"github.com/stackrox/rox/migrator/types"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigration(t *testing.T) {
	cases := []struct {
		image               *storage.Image
		listImage           *storage.ListImage
		expectedCVEs        int32
		expectedFixableCVEs int32
	}{
		{
			image: &storage.Image{
				Id: "image-1",
				Scan: &storage.ImageScan{
					Components: []*storage.EmbeddedImageScanComponent{
						{
							Vulns: []*storage.EmbeddedVulnerability{
								{
									Cve: "cve-1",
								},
								{
									Cve: "cve-2",
									SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{
										FixedBy: "v2",
									},
								},
							},
						},
						{
							Vulns: []*storage.EmbeddedVulnerability{
								{
									Cve: "cve-3",
								},
								{
									Cve: "cve-4",
								},
							},
						},
					},
				},
			},
			listImage: &storage.ListImage{
				Id: "image-1",
			},
			expectedCVEs:        4,
			expectedFixableCVEs: 1,
		},
		{
			image: &storage.Image{
				Id: "image-2",
				Scan: &storage.ImageScan{
					Components: []*storage.EmbeddedImageScanComponent{
						{
							Vulns: []*storage.EmbeddedVulnerability{
								{
									Cve: "cve-1",
								},
								{
									Cve: "cve-2",
									SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{
										FixedBy: "v2",
									},
								},
							},
						},
						{
							Vulns: []*storage.EmbeddedVulnerability{
								{
									Cve: "cve-1",
								},
								{
									Cve: "cve-2",
								},
							},
						},
					},
				},
			},
			listImage: &storage.ListImage{
				Id: "image-2",
			},
			expectedCVEs:        2,
			expectedFixableCVEs: 1,
		},
	}

	db, err := bolthelpers.NewTemp(testutils.DBFileNameForT(t))
	require.NoError(t, err)

	badgerDB, err := badgerhelpers.NewTemp("temp")
	require.NoError(t, err)

	err = fillImages(badgerDB, []*storage.Image{cases[0].image, cases[1].image})
	require.NoError(t, err)

	require.NoError(t, rewriteImagesWithCorrectScanStats(&types.Databases{BoltDB: db, BadgerDB: badgerDB}))

	for _, c := range cases {
		validateMigration(t, badgerDB, c.image, c.expectedCVEs, c.expectedFixableCVEs)
	}
}

func fillImages(db *badger.DB, images []*storage.Image) error {
	for _, image := range images {
		err := db.Update(func(tx *badger.Txn) error {
			key := getKey(imageBucketName, image.GetId())

			data, err := proto.Marshal(image)
			if err != nil {
				return err
			}
			return tx.Set(key, data)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func validateMigration(t *testing.T, db *badger.DB, image *storage.Image, expectedCVEs, expectedFixableCVEs int32) {
	err := db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(getKey(imageBucketName, image.GetId()))
		require.NoError(t, err)

		image = &storage.Image{}
		err = item.Value(func(v []byte) error {
			return proto.Unmarshal(v, image)
		})
		require.NoError(t, err)

		assert.Equal(t, expectedCVEs, image.GetCves())
		assert.Equal(t, expectedFixableCVEs, image.GetFixableCves())

		item, err = tx.Get(getKey(listImageBucketName, image.GetId()))
		require.NoError(t, err)

		listImage := &storage.ListImage{}
		err = item.Value(func(v []byte) error {
			return proto.Unmarshal(v, listImage)
		})
		require.NoError(t, err)

		assert.Equal(t, expectedCVEs, listImage.GetCves())
		assert.Equal(t, expectedFixableCVEs, listImage.GetFixableCves())
		return nil
	})
	require.NoError(t, err)
}