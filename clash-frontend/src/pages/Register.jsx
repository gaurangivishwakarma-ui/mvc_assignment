import SystemPopup from '../components/SystemPopup';
import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { registerPlayer } from '../api/auth';

export default function Register() {
    const [username, setUsername] = useState('');
    const [systemPopupMsg, setSystemPopupMsg] = useState(null);
    const showPopup = (msg) => setSystemPopupMsg(msg);
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        setIsLoading(true);
        try {
            await registerPlayer(username, password);
            showPopup("Commander registered successfully! Please log in.");
            navigate('/login');
        } catch (err) {
            setError(err);
            setIsLoading(false);
        }
    };

    const styles = {
        container: {
            backgroundImage: 'url("/background.jpg")',
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            backgroundRepeat: 'no-repeat',
            minHeight: '100vh',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            padding: '20px',
        },
        card: {
            background: 'linear-gradient(to bottom, #e2d1b3, #c5ad82)',
            border: '6px solid #4a3525',
            borderRadius: '12px',
            padding: '40px 30px',
            maxWidth: '420px',
            width: '100%',
            textAlign: 'center',
            boxShadow: 'inset 0px 0px 15px rgba(0,0,0,0.3), 0px 15px 30px rgba(0,0,0,0.8)',
            position: 'relative',
        },
        headerStrap: {
            background: '#4a3525',
            color: '#ffce00',
            position: 'absolute',
            top: '-20px',
            left: '50%',
            transform: 'translateX(-50%)',
            padding: '10px 40px',
            border: '3px solid #2d1f14',
            borderRadius: '8px',
            fontSize: '24px',
            fontWeight: '900',
            textShadow: '2px 2px 0px #000',
            boxShadow: '0px 5px 10px rgba(0,0,0,0.5)',
            whiteSpace: 'nowrap',
        },
        errorBox: {
            background: '#ff4d4d',
            color: '#fff',
            border: '2px solid #990000',
            padding: '12px',
            borderRadius: '6px',
            marginBottom: '20px',
            marginTop: '10px',
            fontSize: '15px',
            fontWeight: '800',
            textShadow: '1px 1px 0px #000',
            boxShadow: '0px 4px 0px #990000',
        },
        linkText: {
            marginTop: '25px',
            fontSize: '16px',
            color: '#4a3525',
            fontWeight: '800',
        },
        link: {
            color: '#2d7acc',
            textDecoration: 'none',
            textShadow: 'none',
            marginLeft: '5px',
        }
    };

    return (
        <div style={styles.container}>
            <SystemPopup message={systemPopupMsg} onClose={() => setSystemPopupMsg(null)} />
            <div style={styles.card}>
                <div style={styles.headerStrap}>NEW COMMANDER</div>

                {error && <div style={styles.errorBox}>⚔️ {error}</div>}

                <form onSubmit={handleSubmit} style={{ marginTop: '20px' }}>
                    <input
                        type="text"
                        placeholder="Choose Commander Name"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        className="coc-input"
                        required
                    />
                    <input
                        type="password"
                        placeholder="Create Secret Code"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className="coc-input"
                        required
                    />
                    <button
                        type="submit"
                        className="coc-button"
                        disabled={isLoading}
                    >
                        {isLoading ? 'FORGING ACCOUNT...' : 'CREATE ACCOUNT'}
                    </button>
                </form>
                <div style={styles.linkText}>
                    Already have a base? <Link to="/login" style={styles.link}>Login Here</Link>
                </div>
            </div>
        </div>
    );
}