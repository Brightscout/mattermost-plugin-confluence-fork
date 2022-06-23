import React, {useEffect, useState, useCallback} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {Modal, Button} from 'react-bootstrap';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import selectors from '../../selectors';
import ConfluenceInstanceSelector from '../confluence_instance_selector';
import ConfluenceSpaceSelector from '../confluence_space_selector';
import Validator from '../validator';
import ConfluenceField from '../confluence_field';
import {getModalStyles} from 'src/utils/styles';
import {getSpacesForConfluenceURL, createPageForConfluence, closeCreateConfluencePageModal} from 'src/actions';

const getStyle = () => ({
    typeFormControl: {
        resize: 'none',
        height: '10em',
    },
});

const CreateConfluencePage = (theme) => {
    const dispatch = useDispatch();
    const validator = new Validator();

    const postMessage = useSelector((state) => selectors.isCreateConfluencePageModalVisible(state));
    const channelID = useSelector((state) => getCurrentChannelId(state));

    const [modalVisible, setModalVisible] = useState(false);
    const [instanceID, setInstanceID] = useState('');
    const [pageTitle, setPageTitle] = useState('');
    const [pageDescription, setPageDescription] = useState(postMessage.message);
    const [spaceKey, setSpaceKey] = useState('');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState(false);

    useEffect(() => {
        if (postMessage && postMessage.message) {
            setModalVisible(true);
            setPageDescription(postMessage.message);
        } else {
            setModalVisible(false);
        }
    }, [postMessage]);

    useEffect(() => {
        if (!instanceID) {
            return;
        }
        let response;
        (async () => {
            response = await getSpacesForConfluenceURL(instanceID)(dispatch);
            if (response?.error) {
                setError(response.error.response?.text);
            }
        })();
        setSpaceKey('');
    }, [instanceID]);

    const reset = () => {
        setSaving(false);
        setInstanceID('');
        setSpaceKey('');
        setPageTitle('');
        setPageDescription('');
        setError('');
    };

    const handleClose = useCallback((e) => {
        if (e?.preventDefault) {
            e.preventDefault();
        }
        reset();
        dispatch(closeCreateConfluencePageModal());
      }, [])

    const handleInstanceChange = useCallback((currentInstanceID) => {
        setInstanceID(currentInstanceID);
        setSpaceKey('');
        setError('');
      }, [instanceID])

    const handleSpaceKeyChange = useCallback((currentSpaceKey) => {
        setSpaceKey(currentSpaceKey);
      }, [spaceKey])

    const handlePageTitle = useCallback((e) => {
        setPageTitle(e.target.value);
      }, [pageTitle])

    const handlePageDescription = useCallback((e) => {
        setPageDescription(e.target.value);
      }, [pageDescription])

    const handleSubmit = () => {
        if (!validator.validate()) {
            return;
        }

        const pageDetails = {
            title: pageTitle,
            description: pageDescription,
        };

        setSaving(true);
        (async () => {
            const response = await createPageForConfluence(instanceID, channelID, spaceKey, pageDetails)(dispatch);
            if (response?.error) {
                setError(response.error.response?.text);
                setSaving(false);
            } else {
                reset();
                dispatch(closeCreateConfluencePageModal());
            }
        })();
    };

    return (
        <Modal
            dialogClassName='modal--scroll'
            show={modalVisible}
            onHide={handleClose}
            onExited={handleClose}
            backdrop={'static'}
            bsSize='large'
        >
            <Modal.Header closeButton={true}>
                <Modal.Title>
                    {'Create Confluence Page'}
                </Modal.Title>
            </Modal.Header>
            <Modal.Body style={getModalStyles.modalBody}>
                <ConfluenceInstanceSelector
                    theme={theme}
                    selectedInstanceID={instanceID}
                    onInstanceChange={handleInstanceChange}
                />

                {instanceID && <ConfluenceSpaceSelector
                    theme={theme}
                    selectedSpaceKey={spaceKey}
                    onSpaceKeyChange={handleSpaceKeyChange}
                />}

                {spaceKey && <ConfluenceField
                    label={'Page Title'}
                    type={'text'}
                    fieldType={'input'}
                    required={true}
                    placeholder={'Enter Page Title.'}
                    value={pageTitle}
                    addValidation={validator.addValidation}
                    removeValidation={validator.removeValidation}
                    onChange={handlePageTitle}
                />}
                {spaceKey &&<ConfluenceField
                    label={'Page Description'}
                    formControlStyle={getStyle().typeFormControl}
                    type={'textarea'}
                    fieldType={'input'}
                    required={true}
                    value={pageDescription}
                    addValidation={validator.addValidation}
                    removeValidation={validator.removeValidation}
                    onChange={handlePageDescription}
                />}
                {error && <p className='alert alert-danger'>
                    <i
                        className='fa fa-warning'
                        title='Warning Icon'
                    />
                    <span> {error}</span>
                </p>}
            </Modal.Body>

            {spaceKey && <Modal.Footer style={getModalStyles.modalFooter}>
                <Button
                    type='button'
                    bsStyle='link'
                    onClick={handleClose}
                >
                    {'Cancel'}
                </Button>
                <Button
                    type='submit'
                    bsStyle='primary'
                    onClick={handleSubmit}
                    disabled={saving}
                >
                    {saving && <span className='fa fa-spinner fa-fw fa-pulse spinner'/>}
                    {'Save Subscription'}
                </Button>
            </Modal.Footer>}
        </Modal>
    );
};

export default CreateConfluencePage;
