import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';
import { reduxForm, formValueSelector, change } from 'redux-form';
import sortBy from 'lodash/sortBy';

import Select from 'Components/ReactSelect';
import flattenObject from 'utils/flattenObject';
import { removeEmptyPolicyFields } from 'utils/policyUtils';
import { getPolicyFormDataKeys } from 'Containers/Policies/Wizard/Form/utils';
import FormField from 'Components/FormField';
import Field from 'Containers/Policies/Wizard/Form/Field';

class FieldGroupCards extends Component {
    static propTypes = {
        fieldGroups: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
        formData: PropTypes.shape({}).isRequired,
        change: PropTypes.func.isRequired,
    };

    constructor(props) {
        super(props);

        this.state = {
            fields: [],
        };
    }

    addFormField = (option) => {
        let fieldToAdd = {};

        this.props.fieldGroups.forEach((fieldGroup) => {
            const field = fieldGroup.descriptor.find((obj) => obj.jsonpath === option);
            if (field) {
                fieldToAdd = field;
            }
        });
        this.setState((prevState) => ({ fields: prevState.fields.concat(fieldToAdd.jsonpath) }));
        if (fieldToAdd.defaultValue !== undefined && fieldToAdd.defaultValue !== null) {
            this.props.change(fieldToAdd.jsonpath, fieldToAdd.defaultValue);
        }
    };

    removeField = (jsonpath) => {
        let fieldToRemove = {};
        this.props.fieldGroups.forEach((fieldGroup) => {
            const field = fieldGroup.descriptor.find((obj) => obj.jsonpath === jsonpath);

            if (field) {
                fieldToRemove = field;
            }
        });

        this.setState((prevState) => ({
            fields: prevState.fields.filter((fieldPath) => fieldPath !== fieldToRemove.jsonpath),
        }));

        this.props.change(fieldToRemove.jsonpath, null);
    };

    renderFields = (formFields, formData) => {
        const filteredFields = formFields.filter((field) => {
            const isAddedField =
                this.state.fields.length !== 0 &&
                this.state.fields.find((o) => o === field.jsonpath);
            return (
                !field.header &&
                (field.default ||
                    isAddedField ||
                    formData.find((jsonpath) => jsonpath.includes(field.jsonpath)))
            );
        });

        if (!filteredFields.length) {
            if (this.isHeaderOnlyCard(formFields)) {
                return '';
            }

            return <div className="p-3 text-base-500 font-500">No Fields Added</div>;
        }

        return (
            <div className="h-full p-3">
                {filteredFields.map((field) => {
                    const removeField = !field.default ? this.removeField : null;
                    return (
                        <FormField
                            key={field.jsonpath}
                            label={field.label}
                            name={field.jsonpath}
                            required={field.required}
                            onRemove={removeField}
                        >
                            <Field field={field} />
                        </FormField>
                    );
                })}
            </div>
        );
    };

    renderFieldsDropdown = (formFields, formData) => {
        let availableFields = formFields
            .filter(
                (field) =>
                    !this.state.fields.find((jsonpath) => jsonpath === field.jsonpath) &&
                    !field.default &&
                    !formData.find((jsonpath) => jsonpath.includes(field.jsonpath))
            )
            .map((field) => ({ label: field.label, value: field.jsonpath }));
        const placeholder = 'Add a field';
        if (!availableFields.length) {
            return '';
        }
        availableFields = sortBy(availableFields, (o) => o.label);
        return (
            <div className="flex p-3 border-t border-base-200 bg-success-100">
                <span className="w-full">
                    <Select
                        id="policyConfigurationSelect"
                        onChange={this.addFormField}
                        options={availableFields}
                        placeholder={placeholder}
                        menuPlacement="auto"
                    />
                </span>
            </div>
        );
    };

    renderHeaderControl = (formFields) => {
        const headerField = formFields.find((field) => field.header);
        if (!headerField) {
            return '';
        }

        return (
            <div className="header-control float-right">
                {headerField.label && (
                    <label htmlFor={headerField.jsonpath} className="pr-1">
                        {headerField.label}
                    </label>
                )}
                <Field field={headerField} />
            </div>
        );
    };

    isHeaderOnlyCard = (formFields) =>
        formFields.length === 1 && formFields.find((field) => field.header);

    render() {
        const formData = Object.keys(flattenObject(removeEmptyPolicyFields(this.props.formData)));

        return this.props.fieldGroups.map((fieldGroup) => {
            const { header: fieldGroupName, descriptor: formFields, dataTestId } = fieldGroup;
            const headerControl = this.renderHeaderControl(formFields);
            const border = this.isHeaderOnlyCard(formFields) ? '' : 'border';
            const leading = headerControl ? 'leading-loose' : 'leading-normal';
            return (
                <div className="px-3 pt-5" data-testid={dataTestId} key={fieldGroupName}>
                    <div className={`bg-base-100 ${border} border-base-200 shadow`}>
                        <div
                            className={`p-2 pb-2 border-b border-base-300 text-base-600 font-700 text-lg capitalize flex justify-between ${leading}`}
                        >
                            {fieldGroupName}
                            {headerControl}
                        </div>
                        {this.renderFields(formFields, formData)}
                        {this.renderFieldsDropdown(formFields, formData)}
                    </div>
                </div>
            );
        });
    }
}

const formFields = (state) =>
    formValueSelector('policyCreationForm')(state, ...getPolicyFormDataKeys());

const getFormData = createSelector([formFields], (formData) => formData);

const mapStateToProps = createStructuredSelector({
    formData: getFormData,
});

const mapDispatchToProps = (dispatch) => ({
    change: (field, value) => dispatch(change('policyCreationForm', field, value)),
});

export default reduxForm({
    form: 'policyCreationForm',
    enableReinitialize: true,
    destroyOnUnmount: false,
    keepDirtyOnReinitialize: true,
})(connect(mapStateToProps, mapDispatchToProps)(FieldGroupCards));