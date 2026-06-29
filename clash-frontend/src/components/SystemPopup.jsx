import React from 'react';

export default function SystemPopup({ message, onClose }) {
    if (!message) return null;

    return (
        <div style={{
            position: 'fixed',
            top: 0, left: 0, right: 0, bottom: 0,
            backgroundColor: 'rgba(0,0,0,0.85)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            zIndex: 99999
        }}>
            <div style={{
                background: 'linear-gradient(to bottom, #4a2511, #2a1105)',
                border: '4px solid #ffce00',
                borderRadius: '12px',
                padding: '30px',
                width: '350px',
                textAlign: 'center',
                boxShadow: '0px 15px 40px rgba(0,0,0,0.9), inset 0px 0px 20px rgba(0,0,0,0.5)',
                color: '#fff',
                fontFamily: "'Montserrat', sans-serif"
            }}>
                <div style={{
                    fontSize: '22px',
                    fontWeight: '900',
                    color: '#ffce00',
                    marginBottom: '20px',
                    textTransform: 'uppercase',
                    textShadow: '2px 2px 0px #000'
                }}>
                    Notice
                </div>
                <div style={{
                    fontSize: '16px',
                    lineHeight: '1.5',
                    marginBottom: '25px',
                    color: '#f0e6d2'
                }}>
                    {message}
                </div>
                <button 
                    className="coc-button"
                    style={{
                        background: 'linear-gradient(to bottom, #f29f05, #d96d00)',
                        border: '2px solid #5a2a00',
                        padding: '10px 30px',
                        fontSize: '18px',
                        fontWeight: 'bold',
                        color: '#fff',
                        borderRadius: '6px',
                        cursor: 'pointer',
                        textShadow: '1px 1px 0px #000',
                        boxShadow: 'inset 0px 2px 0px rgba(255,255,255,0.4), 0px 4px 0px #5a2a00'
                    }}
                    onClick={onClose}
                >
                    OKAY
                </button>
            </div>
        </div>
    );
}
