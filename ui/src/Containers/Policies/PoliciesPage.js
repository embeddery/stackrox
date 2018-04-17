import React, { Component } from 'react';
import ReactRouterPropTypes from 'react-router-prop-types';
import { ToastContainer, toast } from 'react-toastify';
import { Form } from 'react-form';
import axios from 'axios';
import * as Icon from 'react-feather';
import isEqual from 'lodash/isEqual';
import find from 'lodash/find';
import Table from 'Components/Table';
import Panel from 'Components/Panel';
import PolicyCreationForm from 'Containers/Policies/PolicyCreationForm';
import PolicyView from 'Containers/Policies/PoliciesView';
import PoliciesPreview from 'Containers/Policies/PoliciesPreview';

import { severityLabels } from 'messages/common';
import { sortSeverity } from 'sorters/sorters';

const reducer = (action, prevState, nextState) => {
    switch (action) {
        case 'UPDATE_POLICIES':
            return { policies: nextState.policies };
        case 'SELECT_POLICY':
            return {
                selectedPolicy: nextState.policy,
                editingPolicy: null,
                addingPolicy: false,
                showPreviewPolicy: false
            };
        case 'UNSELECT_POLICY':
            return { selectedPolicy: null, editingPolicy: null, addingPolicy: false };
        case 'EDIT_POLICY':
            return { editingPolicy: nextState.policy };
        case 'ADD_POLICY':
            return { editingPolicy: nextState.policy, addingPolicy: true };
        case 'CANCEL_EDIT_POLICY':
            return { editingPolicy: null, addingPolicy: false };
        case 'CANCEL_ADD_POLICY':
            return { selectedPolicy: null, editingPolicy: null, addingPolicy: false };
        default:
            return prevState;
    }
};

class PoliciesPage extends Component {
    static propTypes = {
        match: ReactRouterPropTypes.match.isRequired
    };

    constructor(props) {
        super(props);

        this.state = {
            policies: [],
            policyCategories: [],
            notifiers: [],
            clusters: [],
            deployments: [],
            selectedPolicy: null,
            editingPolicy: null,
            addingPolicy: false,
            policyDryRun: null,
            showPreviewPolicy: false
        };
    }

    componentDidMount() {
        this.pollPolicies();
        this.retrievePolicyCategories();
        this.retrieveNotifiers();
        this.retrieveClusters();
        this.retrieveDeployments();
    }

    componentWillUnmount() {
        if (this.pollTimeoutId) {
            clearTimeout(this.pollTimeoutId);
            this.pollTimeoutId = null;
        }
    }

    onSubmit = policy => {
        if (this.state.addingPolicy) this.createPolicy(policy);
        else this.savePolicy(policy);
    };

    getPolicyFromId = policies => {
        if (this.props.match.params.id) {
            policies.map(policy => {
                if (policy.id === this.props.match.params.id) {
                    this.selectPolicy(policy);
                }
                return null;
            });
        }
    };

    getPolicies = () =>
        axios.get('/v1/policies', { params: this.params }).then(response => {
            if (!response.data.policies || isEqual(this.state.policies, response.data.policies))
                return;
            const { policies } = response.data;
            this.update('UPDATE_POLICIES', { policies });
            this.getPolicyFromId(policies);
        });

    getPolicyDryRun = policy => {
        let filteredPolicy = this.policyCreationForm.removeEmptyFields(policy);
        filteredPolicy = this.policyCreationForm.postFormatWhitelistField(filteredPolicy);
        filteredPolicy = this.policyCreationForm.postFormatScopeField(filteredPolicy);
        this.update('EDIT_POLICY', { policy: filteredPolicy });
        axios
            .post('/v1/policies/dryrun', filteredPolicy)
            .then(response => {
                if (!response.data) return;
                const policyDryRun = response.data;
                this.setState({ policyDryRun, showPreviewPolicy: true });
            })
            .catch(error => {
                console.error(error);
                if (error.response) toast(error.response.data.error);
            });
    };

    getPolicyCategories = () => axios.get('/v1/policyCategories');

    getNotifiers = () => axios.get('/v1/notifiers');

    getClusters = () => axios.get('/v1/clusters');

    retrievePolicyCategories = () => {
        this.getPolicyCategories().then(response => {
            if (!response.data.categories) return;
            const { categories } = response.data;
            const policyCategories = categories;
            this.setState({ policyCategories });
        });
    };

    retrieveNotifiers = () => {
        this.getNotifiers().then(response => {
            if (!response.data.notifiers) return;
            const { notifiers } = response.data;
            this.setState({ notifiers });
        });
    };

    retrieveClusters = () => {
        this.getClusters().then(response => {
            if (!response.data.clusters) return;
            const { clusters } = response.data;
            this.setState({ clusters });
        });
    };

    retrieveDeployments = () => {
        axios.get('/v1/deployments').then(response => {
            if (!response.data.deployments) return;
            const { deployments } = response.data;
            this.setState({ deployments });
        });
    };

    pollPolicies = () => {
        this.getPolicies().then(() => {
            this.pollTimeoutId = setTimeout(this.pollPolicies, 5000);
        });
    };

    preSubmit = policy => this.policyCreationForm.preSubmit(policy);

    reassessPolicies = () => {
        axios
            .post('/v1/policies/reassess')
            .then(() => {
                toast('Policies were reassessed');
                this.policyTable.clearSelectedRows();
            })
            .catch(error => {
                console.error(error);
                if (error.response) toast(error.response.data.error);
            });
    };

    deletePolicies = () => {
        const promises = [];
        this.policyTable.getSelectedRows().forEach(obj => {
            // close the view panel if that policy is being deleted
            if (this.state.selectedPolicy && obj.id === this.state.selectedPolicy.id) {
                this.unselectPolicy();
            }
            const promise = axios.delete(`/v1/policies/${obj.id}`);
            promises.push(promise);
        });
        Promise.all(promises).then(() => {
            this.policyTable.clearSelectedRows();
            this.getPolicies();
        });
    };

    addPolicy = () => {
        this.update('ADD_POLICY', { policy: {} });
    };

    selectPolicy = policy => {
        this.update('SELECT_POLICY', { policy });
    };

    unselectPolicy = () => {
        this.update('UNSELECT_POLICY');
    };

    editPolicy = policy => {
        this.update('EDIT_POLICY', { policy });
    };

    cancelAddingPolicy = () => {
        this.update('CANCEL_ADD_POLICY');
    };

    cancelEditingPolicy = () => {
        this.update('CANCEL_EDIT_POLICY');
    };

    createPolicy = policy => {
        axios
            .post('/v1/policies', policy)
            .then(() => {
                this.cancelAddingPolicy();
                this.getPolicies().then(() => {
                    this.selectPolicy(find(this.state.policies, ['name', policy.name]));
                });
            })
            .catch(error => {
                console.error(error);
                if (error.response) toast(error.response.data.error);
            });
    };

    updatePolicy = policy =>
        axios
            .put(`/v1/policies/${policy.id}`, policy)
            .then(() => {
                this.getPolicies();
            })
            .catch(error => {
                console.error(error);
                if (error.response) toast(error.response.data.error);
                return error;
            });

    savePolicy = policy => {
        this.updatePolicy(policy).then(error => {
            if (error !== undefined) return;
            this.cancelEditingPolicy();
            this.selectPolicy(policy);
        });
    };

    toggleEnabledDisabledPolicy = policy => {
        const newPolicy = Object.assign({}, policy);
        newPolicy.disabled = !policy.disabled;
        this.updatePolicy(newPolicy);
    };

    update = (action, nextState) => {
        this.setState(prevState => reducer(action, prevState, nextState));
    };

    closePreviewPanel = () => {
        this.setState({ showPreviewPolicy: false });
        return this.cancelEditingPolicy();
    };

    closeEditPanel = () => {
        const newPolicy = this.state.addingPolicy;
        if (newPolicy) return this.cancelAddingPolicy();
        return this.cancelEditingPolicy();
    };

    renderTablePanel = () => {
        const header = `${this.state.policies.length} Policies`;
        const buttons = [
            {
                renderIcon: () => <Icon.Trash2 className="h-4 w-4" />,
                text: 'Delete Policies',
                className:
                    'flex py-1 px-2 rounded-sm text-danger-600 hover:text-white hover:bg-danger-400 uppercase text-center text-sm items-center ml-2 bg-white border-2 border-danger-400',
                onClick: this.deletePolicies,
                disabled: this.state.editingPolicy !== null
            },
            {
                renderIcon: () => <Icon.FileText className="h-4 w-4" />,
                text: 'Reassess Policies',
                className:
                    'flex py-1 px-2 rounded-sm text-success-600 hover:text-white hover:bg-success-400 uppercase text-center text-sm items-center ml-2 bg-white border-2 border-success-400',
                onClick: this.reassessPolicies,
                disabled: this.state.editingPolicy !== null,
                tooltip: 'Manually enrich external data'
            },
            {
                renderIcon: () => <Icon.Plus className="h-4 w-4" />,
                text: 'Add',
                className:
                    'flex py-1 px-2 rounded-sm text-success-600 hover:text-white hover:bg-success-400 uppercase text-center text-sm items-center ml-2 bg-white border-2 border-success-400',
                onClick: this.addPolicy,
                disabled: this.state.editingPolicy !== null
            }
        ];
        const columns = [
            {
                key: 'name',
                keys: ['name', 'disabled'],
                keyValueFunc: (name, disabled) => (
                    <div className="flex items-center relative">
                        <div
                            className={`h-2 w-2 rounded-lg absolute -ml-4 ${
                                !disabled ? 'bg-success-500' : 'bg-base-300'
                            }`}
                        />
                        <div>{name}</div>
                    </div>
                ),
                label: 'Name'
            },
            { key: 'description', label: 'Description' },
            {
                key: 'severity',
                label: 'Severity',
                keyValueFunc: severity => severityLabels[severity],
                classFunc: severity => {
                    switch (severity) {
                        case 'Low':
                            return 'text-low-500';
                        case 'Medium':
                            return 'text-medium-500';
                        case 'High':
                            return 'text-high-severity';
                        case 'Critical':
                            return 'text-critical-severity';
                        default:
                            return '';
                    }
                },
                sortMethod: sortSeverity
            }
        ];
        const actions = [
            {
                renderIcon: row =>
                    row.disabled ? (
                        <Icon.Power className="h-5 w-4 text-base-600" />
                    ) : (
                        <Icon.Power className="h-5 w-4 text-success-500" />
                    ),
                className: 'flex rounded-sm uppercase text-center text-sm items-center',
                onClick: this.toggleEnabledDisabledPolicy
            }
        ];
        const rows = this.state.policies;
        return (
            <Panel header={header} buttons={buttons}>
                <Table
                    columns={columns}
                    rows={rows}
                    onRowClick={this.selectPolicy}
                    actions={actions}
                    checkboxes
                    ref={table => {
                        this.policyTable = table;
                    }}
                />
            </Panel>
        );
    };

    renderViewPanel = () => {
        const { notifiers } = this.state;
        const policy = this.state.selectedPolicy;
        const hide = this.state.selectedPolicy === null || this.state.editingPolicy !== null;
        if (hide) return '';
        const header = this.state.selectedPolicy.name;
        const buttons = [
            {
                renderIcon: () => <Icon.Edit className="h-4 w-4" />,
                text: 'Edit',
                className:
                    'flex py-1 px-2 rounded-sm text-success-600 hover:text-white hover:bg-success-400 uppercase text-center text-sm items-center ml-2 bg-white border-2 border-success-400',
                onClick: () => {
                    const { selectedPolicy } = this.state;
                    this.editPolicy(selectedPolicy);
                }
            }
        ];
        return (
            <Panel header={header} buttons={buttons} onClose={this.unselectPolicy} width="w-2/3">
                <PolicyView notifiers={notifiers} policy={policy} />
            </Panel>
        );
    };

    renderEditPanel = () => {
        const { notifiers } = this.state;
        const { clusters } = this.state;
        const { deployments } = this.state;
        const { policyCategories } = this.state;
        const policy = this.state.editingPolicy;
        const hide = policy === null;
        if (hide || this.state.showPreviewPolicy) return '';
        const header = this.state.editingPolicy.name;
        const buttons = [
            {
                renderIcon: () => <Icon.ArrowRight className="h-4 w-4" />,
                text: 'Next',
                className:
                    'flex py-1 px-2 rounded-sm text-primary-600 hover:text-white hover:bg-primary-400 uppercase text-center text-sm items-center ml-2 bg-white border-2 border-primary-400',
                onClick: () => {
                    this.getPolicyDryRun(this.formApi.values);
                }
            }
        ];
        return (
            <Panel header={header} buttons={buttons} onClose={this.closeEditPanel} width="w-2/3">
                <Form onSubmit={this.onSubmit} preSubmit={this.preSubmit}>
                    {formApi => (
                        <PolicyCreationForm
                            clusters={clusters}
                            deployments={deployments}
                            notifiers={notifiers}
                            policy={policy}
                            policyCategories={policyCategories}
                            formApi={formApi}
                            ref={policyCreationForm => {
                                this.policyCreationForm = policyCreationForm;
                                this.formApi = formApi;
                            }}
                        />
                    )}
                </Form>
            </Panel>
        );
    };

    renderPreviewPanel = () => {
        if (!this.state.showPreviewPolicy) return '';
        const policy = this.state.editingPolicy;
        const hide = policy === null;
        if (hide) return '';
        const header = this.state.editingPolicy.name;
        const buttons = [
            {
                renderIcon: () => <Icon.ArrowLeft className="h-4 w-4" />,
                text: 'Previous',
                className:
                    'flex py-1 px-2 rounded-sm text-primary-600 hover:text-white hover:bg-primary-400 uppercase text-center text-sm items-center ml-2 bg-white border-2 border-primary-400',
                onClick: () => {
                    this.setState({ showPreviewPolicy: false });
                }
            },
            {
                renderIcon: () => <Icon.Save className="h-4 w-4" />,
                text: 'Save',
                className:
                    'flex py-1 px-2 rounded-sm text-success-600 hover:text-white hover:bg-success-400 uppercase text-center text-sm items-center ml-2 bg-white border-2 border-success-400',
                onClick: () => {
                    this.setState({ showPreviewPolicy: false });
                    this.onSubmit(this.state.editingPolicy);
                }
            }
        ];
        return (
            <Panel header={header} buttons={buttons} onClose={this.closePreviewPanel} width="w-2/3">
                <PoliciesPreview dryrun={this.state.policyDryRun} />
            </Panel>
        );
    };

    render() {
        return (
            <section className="h-full">
                <ToastContainer
                    toastClassName="font-sans text-base-600 text-white font-600 bg-black"
                    hideProgressBar
                    autoClose={3000}
                />
                <div className="flex flex-1 bg-base-100 h-full">
                    <div className="flex flex-row w-full h-full bg-white rounded-sm shadow">
                        {this.renderTablePanel()}
                        {this.renderViewPanel()}
                        {this.renderEditPanel()}
                        {this.renderPreviewPanel()}
                    </div>
                </div>
            </section>
        );
    }
}

export default PoliciesPage;
