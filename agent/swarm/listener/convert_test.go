package listener

import (
	"testing"
	"time"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestAsDeployment(t *testing.T) {
	t.Parallel()

	time100 := time.Unix(100, 0)

	cases := []struct {
		service  swarm.Service
		expected *v1.Deployment
	}{
		{
			service: swarm.Service{
				ID: "fooID",
				Meta: swarm.Meta{
					Version: swarm.Version{
						Index: 100,
					},
				},
				Endpoint: swarm.Endpoint{
					Ports: []swarm.PortConfig{
						{
							Name:       "api",
							TargetPort: 80,
							Protocol:   swarm.PortConfigProtocolTCP,
						},
					},
				},
				Spec: swarm.ServiceSpec{
					Annotations: swarm.Annotations{
						Name: "foo",
					},
					Mode: swarm.ServiceMode{
						Replicated: &swarm.ReplicatedService{
							Replicas: &[]uint64{10}[0],
						},
					},
					TaskTemplate: swarm.TaskSpec{
						ContainerSpec: &swarm.ContainerSpec{
							Args:    []string{"--flags", "--args"},
							Command: []string{"init"},
							Dir:     "/bin",
							Env:     []string{"LOGLEVEL=Warn", "JVMFLAGS=-Xms256m", "invalidENV"},
							Image:   "nginx:latest",
							User:    "root",
							Privileges: &swarm.Privileges{
								SELinuxContext: &swarm.SELinuxContext{
									User:  "user",
									Role:  "role",
									Type:  "type",
									Level: "level",
								},
							},
							Mounts: []mount.Mount{
								{
									Source:   "volumeSource",
									Type:     mount.TypeVolume,
									ReadOnly: true,
									Target:   "/var/data",
								},
							},
						},
					},
				},
				UpdateStatus: &swarm.UpdateStatus{
					CompletedAt: &time100,
				},
			},
			expected: &v1.Deployment{
				Id:        "fooID",
				Name:      "foo",
				Version:   "100",
				Type:      "Replicated",
				Replicas:  10,
				UpdatedAt: &timestamp.Timestamp{Seconds: 100},
				Containers: []*v1.Container{
					{
						Config: &v1.ContainerConfig{
							Args: []string{"--flags", "--args"},
							Env: []*v1.ContainerConfig_EnvironmentConfig{
								{
									Key:   "LOGLEVEL",
									Value: "Warn",
								},
								{
									Key:   "JVMFLAGS",
									Value: "-Xms256m",
								},
							},
							Command:   []string{"init"},
							Directory: "/bin",
							User:      "root",
						},
						Image: &v1.Image{
							Registry: "docker.io",
							Remote:   "library/nginx",
							Tag:      "latest",
						},
						SecurityContext: &v1.SecurityContext{
							Selinux: &v1.SecurityContext_SELinux{
								User:  "user",
								Role:  "role",
								Type:  "type",
								Level: "level",
							},
						},
						Ports: []*v1.PortConfig{
							{
								Name:          "api",
								ContainerPort: 80,
								Protocol:      "tcp",
							},
						},
						Volumes: []*v1.Volume{
							{
								Name:     "volumeSource",
								Type:     "volume",
								ReadOnly: true,
								Path:     "/var/data",
							},
						},
					},
				},
			},
		},
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cases {
		got := serviceWrap(c.service).asDeployment(cli)

		assert.Equal(t, c.expected, got)
	}

}
