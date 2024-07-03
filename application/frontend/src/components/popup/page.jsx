// PopupSection.jsx
import React, { useState } from 'react';
import Modal from 'react-modal';

const PopupSection = ({ isOpen, onRequestClose, title, inputFields, onSubmit }) => {
    const [params, setParams] = useState({});

    const handleChange = (e) => {
        const { name, value } = e.target;
        setParams(prevParams => ({
            ...prevParams,
            [name]: value
        }));
    };

    const handleSubmit = async () => {
        try {
            await onSubmit(params);
            console.log('Action performed successfully');
        } catch (error) {
            console.error('Error performing action:', error);
        }
    };

    return (
        <Modal isOpen={isOpen} onRequestClose={onRequestClose}>
            <h2>{title}</h2>
            <form>
                {inputFields.map(field => (
                    <div key={field.name}>
                        <label>{field.label}:</label>
                        <input type={field.type} name={field.name} value={params[field.name] || ''} onChange={handleChange} />
                    </div>
                ))}
                <button type="button" onClick={handleSubmit}>Perform Action</button>
            </form>
        </Modal>
    );
};

export default PopupSection;
