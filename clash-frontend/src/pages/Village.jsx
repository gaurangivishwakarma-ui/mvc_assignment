import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getVillageProfile, getShopCatalog, getArmyCatalog, getArmyStatus } from '../api/village';


const BUILDING_ASSETS = {
    TOWN_HALL: '/assets/buildings/town_hall.png',
    GOLD_MINE: '/assets/buildings/gold_mine.png',
    ELIXIR_COLLECTOR: '/assets/buildings/elixir_pump.png',
    BARRACKS: '/assets/buildings/barracks.png',
    DEFAULT: '/assets/buildings/construction_placeholder.png'
};

const styles = {
    container: {
        background: '#3f6319',
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'row',
        position: 'relative',
        overflow: 'hidden'
    },
    topBar: {
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '25px 15px',
        background: 'linear-gradient(to right, #4a3525, #2d1f14)',
        borderRight: '4px solid #1a120b',
        boxShadow: '4px 0px 10px rgba(0,0,0,0.5)',
        zIndex: 10
    },
    playerInfo: { color: '#fff', display: 'flex', flexDirection: 'column', gap: '5px', alignItems: 'center', textAlign: 'center' },
    levelBadge: {
        background: '#ffce00',
        color: '#000',
        padding: '5px 10px',
        borderRadius: '12px',
        fontWeight: '900',
        fontSize: '14px',
        border: '2px solid #b38f00',
        display: 'inline-block'
    },
    resources: { display: 'flex', flexDirection: 'column', gap: '20px', alignItems: 'center' },
    resourceBox: {
        background: 'rgba(0,0,0,0.6)',
        border: '2px solid #7c6248',
        borderRadius: '8px',
        padding: '8px 15px',
        color: '#fff',
        fontWeight: '800',
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
        boxShadow: 'inset 0px 2px 4px rgba(0,0,0,0.8)'
    },
    logoutBtn: { width: 'auto', padding: '8px 15px', margin: '0', fontSize: '14px' },
    dashboardBox: {
        background: 'rgba(0,0,0,0.6)',
        border: '2px solid #7c6248',
        borderRadius: '8px',
        padding: '12px',
        color: '#fff',
        fontWeight: '800',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-start',
        gap: '8px',
        marginTop: '15px',
        boxShadow: 'inset 0px 2px 4px rgba(0,0,0,0.8)',
        width: '100%',
        boxSizing: 'border-box'
    },

    mapViewport: {
        flex: 1,
        overflow: 'hidden',
        padding: '20px',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#49731d'
    },
    gridBoard: {
        display: 'grid',

        gridTemplateColumns: 'repeat(25, 1fr)',
        gridTemplateRows: 'repeat(25, 1fr)',
        gap: '1px',
        background: '#558721',


        width: '90vh',
        maxWidth: '100%',
        aspectRatio: '1 / 1',

        borderRadius: '8px',
        border: '6px solid #2e4712',
        boxShadow: '0px 12px 30px rgba(0,0,0,0.4)',
        backgroundImage: 'linear-gradient(rgba(0,0,0,0.06) 1px, transparent 1px), linear-gradient(90deg, rgba(0,0,0,0.06) 1px, transparent 1px)',
        backgroundSize: '4% 4%'
    },
    buildingWrapper: {
        position: 'relative',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        cursor: 'pointer',
        transition: 'transform 0.1s ease'
    },
    buildingSprite: {
        width: '90%',
        height: '90%',
        objectFit: 'contain',
        filter: 'drop-shadow(0px 8px 4px rgba(0,0,0,0.35))'
    },
    buildingLabel: {
        position: 'absolute',
        bottom: '2px',
        background: 'rgba(0,0,0,0.75)',
        color: '#fff',
        padding: '2px 6px',
        borderRadius: '4px',
        fontSize: '10px',
        fontWeight: 'bold',
        pointerEvents: 'none',
        whiteSpace: 'nowrap'
    },

    bottomBar: { display: 'flex', flexDirection: 'column', justifyContent: 'space-between', padding: '20px', background: 'linear-gradient(to left, rgba(0,0,0,0.7), transparent)', borderLeft: '2px solid rgba(255,255,255,0.1)', zIndex: 10 },
    actionButton: { width: '120px', padding: '15px', fontSize: '16px' },
    attackButton: { width: '150px', padding: '15px', fontSize: '20px', background: 'linear-gradient(to bottom, #ff5e5e, #cc0000)', borderColor: '#660000', boxShadow: 'inset 0px 2px 0px rgba(255,255,255,0.4), 0px 6px 0px #4d0000', textShadow: '2px 2px 0px #4d0000' },
    loadingScreen: { flex: 1, display: 'flex', justifyContent: 'center', alignItems: 'center', background: '#315800', color: '#fff', fontSize: '24px', fontWeight: '800' },

    modalOverlay: { position: 'absolute', top: 0, left: 0, right: 0, bottom: 0, backgroundColor: 'rgba(0,0,0,0.75)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 100, padding: '20px' },
    modalBox: { background: 'linear-gradient(to bottom, #3a2512, #24160a)', border: '5px solid #1a1007', borderRadius: '16px', width: '90%', maxWidth: '500px', padding: '20px', boxShadow: '0px 10px 25px rgba(0,0,0,0.8)' },
    modalHeader: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '3px solid #1a1007', paddingBottom: '10px', marginBottom: '25px', color: '#ffce00' },
    shopGrid: { display: 'flex', flexDirection: 'column', gap: '12px' },
    shopItem: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'rgba(0,0,0,0.4)', border: '2px solid #5a3d25', padding: '12px', borderRadius: '8px', color: '#fff' }
};

export default function Village() {
    const navigate = useNavigate();
    const [player, setPlayer] = useState(null);
    const [error, setError] = useState(null);
    const [isShopOpen, setIsShopOpen] = useState(false);
    const [shopCatalog, setShopCatalog] = useState([]);
    const [isArmyOpen, setIsArmyOpen] = useState(false);
    const [armyCatalog, setArmyCatalog] = useState([]);
    const [armyStatus, setArmyStatus] = useState([]);

    useEffect(() => {
        const loadVillage = async () => {
            try {
                const data = await getVillageProfile();
                setPlayer(data);
            } catch (err) {
                console.error("Village Load Error:", err);
                setError("Failed to load village data.");
            }
        };
        const loadShop = async () => {
            try {
                const catalog = await getShopCatalog();
                setShopCatalog(catalog);
            } catch (err) {
                console.error("Shop Load Error:", err);
            }
        };

        const loadArmy = async () => {
            try {
                const catalog = await getArmyCatalog();
                setArmyCatalog(catalog.troops_available || []);
                const status = await getArmyStatus();
                setArmyStatus(status.army || []);
            } catch (err) {
                console.error("Army Load Error:", err);
            }
        };

        loadVillage();
        loadShop();
        loadArmy();
    }, [navigate]);

    const handleBuyBuilding = async (item) => {
        alert(`Purchasing ${item.name}! Processing with coordinates next.`);
        setIsShopOpen(false);
    };

    const handleTrainTroop = async (troop) => {
        alert(`Training ${troop.troop_type}! Level ${troop.level_req}.`);
        setIsArmyOpen(false);
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/login');
    };

    if (error) return <div style={styles.container}> <div style={styles.loadingScreen}>⚔️ {error}</div> </div>;
    if (!player) return <div style={styles.container}> <div style={styles.loadingScreen}>LOADING BASE...</div> </div>;

    const deployedBuildings = player.buildings_placed || [];

    return (
        <div style={styles.container}>
            <div style={styles.topBar}>
                <div style={styles.playerInfo}>
                    <span style={{ fontSize: '20px', fontWeight: '900' }}>{player.username || 'Commander'}</span>
                    <div>
                        <span style={styles.levelBadge}>Town Hall Level {player.village_level || 1}</span>
                    </div>

                    <div style={styles.dashboardBox}>
                        <div style={{ color: '#aaa', fontSize: '12px', textAlign: 'center', width: '100%', borderBottom: '1px solid #444', paddingBottom: '4px', marginBottom: '4px' }}>DASHBOARD</div>
                        <div>🪙 {player.balances?.gold || 0} Gold</div>
                        <div>⚗️ {player.balances?.elixir || 0} Elixir</div>
                        <div>🏠 {deployedBuildings.length} Buildings</div>
                        <div style={{ fontSize: '12px' }}>⚔️ Army: {armyStatus.length > 0 ? armyStatus.map(t => `${t.quantity} ${t.troop_type}`).join(', ') : 'None'}</div>
                    </div>

                    <button className="coc-button" style={{ ...styles.actionButton, ...styles.attackButton, marginTop: '15px' }}>ATTACK!</button>
                </div>

                <div style={styles.resources}>
                    <button className="coc-button" style={styles.logoutBtn} onClick={handleLogout}>LOGOUT</button>
                </div>
            </div>


            <div style={styles.mapViewport}>
                <div style={styles.gridBoard}>
                    {deployedBuildings.map((b, index) => {
                        const imageSrc = BUILDING_ASSETS[b.building_type] || BUILDING_ASSETS.DEFAULT;

                        return (
                            <div
                                key={index}
                                style={{
                                    ...styles.buildingWrapper,
                                    gridColumn: `${b.x_coords || 1} / span ${b.width || 2}`,
                                    gridRow: `${b.y_coords || 1} / span ${b.height || 2}`
                                }}
                            >
                                <img
                                    src={imageSrc}
                                    alt={b.name}
                                    style={styles.buildingSprite}
                                    onError={(e) => {
                                        e.target.style.display = 'none';
                                        e.target.parentNode.style.background = '#8c633e';
                                        e.target.parentNode.style.border = '2px solid gold';
                                    }}
                                />
                                <div style={styles.buildingLabel}>{b.name}</div>
                            </div>
                        );
                    })}
                </div>
            </div>

            <div style={styles.bottomBar}>
                <button className="coc-button" style={styles.actionButton} onClick={() => setIsShopOpen(true)}>SHOP</button>
                <button className="coc-button" style={styles.actionButton} onClick={() => setIsArmyOpen(true)}>ARMY</button>
            </div>

            {isShopOpen && (
                <div style={styles.modalOverlay}>
                    <div style={styles.modalBox}>
                        <div style={styles.modalHeader}>
                            <h2 style={{ margin: 0, fontWeight: '900' }}>VILLAGE SHOP</h2>
                            <button className="coc-button" style={{ width: '40px', margin: 0, padding: '5px' }} onClick={() => setIsShopOpen(false)}>X</button>
                        </div>
                        <div style={styles.shopGrid}>
                            {shopCatalog.map((item, idx) => (
                                <div key={idx} style={styles.shopItem}>
                                    <div>
                                        <div style={{ fontWeight: '900', fontSize: '16px' }}>{item.name}</div>
                                        <div style={{ fontSize: '12px', color: '#aaa' }}>Footprint: {item.width}x{item.breadth} Grid Tiles</div>
                                    </div>
                                    <button
                                        className="coc-button"
                                        style={{ width: 'auto', padding: '8px 15px', fontSize: '13px', margin: 0 }}
                                        onClick={() => handleBuyBuilding(item)}
                                    >
                                        Buy: {item.cost_type === 'elixir' ? '⚗️' : '🪙'} {item.cost}
                                    </button>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            )}

            {isArmyOpen && (
                <div style={styles.modalOverlay}>
                    <div style={styles.modalBox}>
                        <div style={styles.modalHeader}>
                            <h2 style={{ margin: 0, fontWeight: '900' }}>TRAIN TROOPS</h2>
                            <button className="coc-button" style={{ width: '40px', margin: 0, padding: '5px' }} onClick={() => setIsArmyOpen(false)}>X</button>
                        </div>
                        <div style={styles.shopGrid}>
                            {armyCatalog.filter(troop => troop.level_req <= (player?.village_level || 1)).map((troop, idx) => (
                                <div key={idx} style={styles.shopItem}>
                                    <div>
                                        <div style={{ fontWeight: '900', fontSize: '16px', textTransform: 'capitalize' }}>{troop.troop_type} (Lv {troop.level_req})</div>
                                        <div style={{ fontSize: '12px', color: '#aaa' }}>Space: {troop.housing_space} | HP: {troop.hit_points} | Dmg: {troop.damage}</div>
                                    </div>
                                    <button
                                        className="coc-button"
                                        style={{ width: 'auto', padding: '8px 15px', fontSize: '13px', margin: 0 }}
                                        onClick={() => handleTrainTroop(troop)}
                                    >
                                        Train: ⚗️ {troop.elixir_cost}
                                    </button>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}