import SystemPopup from './SystemPopup';
import { useState, useEffect } from 'react';
import { attackOpponent, getArmyStatus } from '../api/village';
import explosionGif from '../assets/explosion.gif';
import swordImg from '../assets/sword.png';
import battleBg from '../assets/battle_bg.png';

const ANIM_MARCH_MS = 1800;
const ANIM_EXPLODE_MS = 1500;
const TOTAL_ANIM_MS = ANIM_MARCH_MS + ANIM_EXPLODE_MS;

export default function BattleOverlay({ matchData, onClose, BUILDING_ASSETS, SHOP_ASSETS, TROOP_ASSETS }) {
    const [army, setArmy] = useState([]);
    const [systemPopupMsg, setSystemPopupMsg] = useState(null);
    const showPopup = (msg) => setSystemPopupMsg(msg);
    const [deployments, setDeployments] = useState({});
    const [selectedTroopKey, setSelectedTroopKey] = useState(null);
    const [gridDeployedTroops, setGridDeployedTroops] = useState([]);


    const [buildingStates, setBuildingStates] = useState({});

    const [battlePhase, setBattlePhase] = useState('deploying');
    const [battleResult, setBattleResult] = useState(null);

    useEffect(() => {
        getArmyStatus()
            .then(s => {
                const available = (s.army || []).filter(t => t.quantity > 0);
                setArmy(available);
                const init = {};
                available.forEach(t => { init[`${t.troop_type}_${t.current_level}`] = t.quantity; });
                setDeployments(init);
                if (available.length > 0) {
                    setSelectedTroopKey(`${available[0].troop_type}_${available[0].current_level}`);
                }
            })
            .catch(err => console.error('Army fetch failed', err));
    }, []);

    const handleGridClick = (e) => {
        if (battlePhase !== 'deploying') return;
        if (!selectedTroopKey) {
            showPopup('Select a troop from the bottom panel first!');
            return;
        }
        const remaining = deployments[selectedTroopKey] || 0;
        if (remaining <= 0) return;

        const rect = e.currentTarget.getBoundingClientRect();
        const xPct = ((e.clientX - rect.left) / rect.width) * 100;
        const yPct = ((e.clientY - rect.top) / rect.height) * 100;

        const [troopType, levelStr] = selectedTroopKey.split('_');

        setGridDeployedTroops(prev => [...prev, {
            id: `${troopType}_${Date.now()}_${prev.length}`,
            troop_type: troopType,
            level: parseInt(levelStr),
            xPct,
            yPct
        }]);

        setDeployments(prev => ({ ...prev, [selectedTroopKey]: prev[selectedTroopKey] - 1 }));
    };

    const handleAttack = async () => {
        if (!gridDeployedTroops.length) { showPopup('Deploy at least 1 troop!'); return; }

        const troopMap = {};
        gridDeployedTroops.forEach(t => {
            const key = `${t.troop_type}_${t.level}`;
            troopMap[key] = (troopMap[key] || 0) + 1;
        });
        const troops = Object.keys(troopMap).map(k => {
            const [ttype, lvlStr] = k.split('_');
            return { troop_type: ttype, level: parseInt(lvlStr), quantity: troopMap[k] };
        });

        setBattlePhase('simulating');

        try {
            const result = await attackOpponent({ opponent_id: matchData.opponent_id, deployed_troops: troops });

            const townHall = {
                building_type: 'town_hall',
                x_coords: 11, y_coords: 11, width: 4, breadth: 4,
            };
            const layout = [townHall, ...(matchData.village_layout || [])];
            const totalBuildings = layout.length;
            const numDestroyed = Math.floor(totalBuildings * result.destruction_percent / 100);

            let indices = Array.from({ length: totalBuildings - 1 }, (_, i) => i + 1);
            indices.sort(() => Math.random() - 0.5);
            const destroyedIndices = indices.slice(0, numDestroyed);

            setGridDeployedTroops(prev => prev.map((t, i) => {
                const targetIdx = destroyedIndices.length > 0 ? destroyedIndices[i % destroyedIndices.length] : 0;
                const b = layout[targetIdx];
                const targetXPct = ((b.x_coords + ((b.width || 2) / 2)) / 25) * 100;
                const targetYPct = ((b.y_coords + ((b.breadth || b.height || 2) / 2)) / 25) * 100;
                return { ...t, targetXPct, targetYPct };
            }));


            setTimeout(() => {
                const explosionDuration = 3500;

                destroyedIndices.forEach((bIdx, i) => {
                    const delay = numDestroyed > 0 ? (explosionDuration / numDestroyed) * i : 0;

                    setTimeout(() => {
                        setBuildingStates(prev => ({ ...prev, [bIdx]: 'exploding' }));

                        setTimeout(() => {
                            setBuildingStates(prev => ({ ...prev, [bIdx]: 'rubble' }));
                        }, 1000);
                    }, delay);
                });

                setTimeout(() => {
                    setBattlePhase('finished');
                    setBattleResult(result);
                }, explosionDuration + 1200);

            }, 1800);

        } catch (err) {
            showPopup(`Attack Failed: ${err}`);
            setBattlePhase('deploying');
        }
    };

    const renderBuildings = () => {
        const tiles = [];
        for (let i = 0; i < 625; i++) {
            const x = (i % 25) + 1;
            const y = Math.floor(i / 25) + 1;
            tiles.push(
                <div key={`tile-${i}`} style={{
                    gridColumn: x, gridRow: y,
                    backgroundImage: 'url(/assets/grass_tile.png)',
                    backgroundSize: 'cover', zIndex: 1
                }} />
            );
        }

        const townHall = {
            building_type: 'town_hall',
            name: `Town Hall Lv${matchData.village_level || 1}`,
            x_coords: 11, y_coords: 11, width: 4, breadth: 4,
            current_level: matchData.village_level || 1
        };
        [townHall, ...(matchData.village_layout || [])].forEach((b, i) => {
            const lvl = b.current_level || 1;
            const type = (b.building_type || '').toLowerCase();
            const key = lvl > 1 ? `${type}_${lvl}` : type;
            const src = SHOP_ASSETS[key] || SHOP_ASSETS[type] || BUILDING_ASSETS[type.toUpperCase()] || BUILDING_ASSETS.DEFAULT;
            const w = b.width || 2;
            const h = b.breadth || b.height || 2;
            const isHall = type === 'town_hall';

            const bState = buildingStates[i] || 'normal';


            let imgFilter = 'drop-shadow(0px 4px 4px rgba(0,0,0,0.5))';
            if (bState === 'exploding') {
                imgFilter = 'brightness(0.35) sepia(1) hue-rotate(-30deg) saturate(5)';
            } else if (bState === 'rubble') {
                imgFilter = 'grayscale(1) brightness(0.25) contrast(1.5) sepia(0.3)';
            }


            const isTownHallExploding = isHall && bState === 'exploding';

            tiles.push(
                <div key={`b-${i}`} style={{
                    position: 'relative',
                    gridColumn: `${b.x_coords || 1} / span ${w}`,
                    gridRow: `${b.y_coords || 1} / span ${h}`,
                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                    zIndex: 3,
                    animation: isTownHallExploding ? 'shake 0.4s infinite' : 'none'
                }}>
                    <img
                        src={src}
                        alt={type}
                        style={{
                            width: '90%', height: '90%', objectFit: 'contain',
                            filter: imgFilter,
                            transition: 'filter 0.3s'
                        }}
                        onError={e => { e.target.style.display = 'none'; e.target.parentNode.style.background = '#8c633e'; }}
                    />

                    {bState === 'exploding' && (
                        <img src={explosionGif} alt="boom" key={`exp-${Date.now()}`}
                            style={{
                                position: 'absolute', width: isHall ? '260%' : '180%', height: isHall ? '260%' : '180%',
                                top: '50%', left: '50%', transform: 'translate(-50%,-50%)',
                                zIndex: 10, pointerEvents: 'none', mixBlendMode: 'screen', opacity: isHall ? 1 : 0.8
                            }}
                        />
                    )}
                </div>
            );
        });

        return tiles;
    };

    return (
        <div style={{
            position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh',
            zIndex: 500,
            backgroundImage: `url(${battleBg})`,
            backgroundSize: 'cover', backgroundPosition: 'center',
            display: 'flex', flexDirection: 'column',
            overflow: 'hidden'
        }}>
            {battlePhase !== 'finished' && (
                <div style={{
                    width: '100%', padding: '10px 0', textAlign: 'center', flexShrink: 0,
                    background: 'linear-gradient(90deg, #8b0000, #cc2200, #ff5500, #cc2200, #8b0000)',
                    backgroundSize: '300% 100%',
                    boxShadow: '0 4px 20px rgba(255,50,0,0.7)',
                    animation: 'shimmer 2s linear infinite, bannerPulse 0.5s ease-in-out infinite alternate',
                    zIndex: 600
                }}>
                    <span style={{
                        fontSize: '20px', fontWeight: '900', color: '#fff',
                        textShadow: '0 0 12px #ffce00, 2px 2px 0 #4d0000',
                        letterSpacing: '4px'
                    }}>
                        ⚔️ &nbsp; {battlePhase === 'simulating' ? 'BATTLE IN PROGRESS' : 'TAP TO DEPLOY'} &nbsp; ⚔️
                    </span>
                </div>
            )}

            <div style={{
                display: 'flex', justifyContent: 'space-between', padding: '12px 30px',
                background: 'rgba(0,0,0,0.65)', color: '#fff', alignItems: 'center',
                boxShadow: '0 4px 10px rgba(0,0,0,0.5)', flexShrink: 0
            }}>
                <div>
                    <h2 style={{ margin: 0, color: '#ffce00', fontSize: '20px' }}>{matchData.username}'s Village</h2>
                    <span style={{ color: '#bbb', fontSize: '13px' }}>Level {matchData.village_level} | XP {matchData.xp_points}</span>
                </div>
                <div style={{ display: 'flex', gap: '12px', alignItems: 'center', fontWeight: 'bold' }}>
                    <div style={{ background: 'rgba(0,0,0,0.5)', padding: '5px 14px', borderRadius: '20px', border: '1px solid gold' }}>
                        🪙 {matchData.loot_available?.gold_coins || 0}
                    </div>
                    <div style={{ background: 'rgba(0,0,0,0.5)', padding: '5px 14px', borderRadius: '20px', border: '1px solid #d85bd8' }}>
                        ⚗️ {matchData.loot_available?.elixir || 0}
                    </div>
                    {battlePhase === 'deploying' && (
                        <button className="coc-button" onClick={onClose}
                            style={{ padding: '10px 20px', margin: 0, fontSize: '14px' }}>
                            SURRENDER
                        </button>
                    )}
                </div>
            </div>

            <div style={{ flex: 1, position: 'relative', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <div
                    onClick={handleGridClick}
                    style={{
                        position: 'relative',
                        display: 'grid',
                        gridTemplateColumns: 'repeat(25, 1fr)',
                        gridTemplateRows: 'repeat(25, 1fr)',
                        width: 'min(72vh, 80vw)', height: 'min(72vh, 80vw)',
                        flexShrink: 0,
                        border: '3px solid rgba(30,60,10,0.6)',
                        borderRadius: '6px',
                        background: 'rgba(0,0,0,0.08)',
                        boxShadow: '0 0 40px rgba(0,0,0,0.5)',
                        cursor: battlePhase === 'deploying' ? 'crosshair' : 'default'
                    }}
                >
                    {renderBuildings()}


                    {gridDeployedTroops.map((t) => {
                        const src = TROOP_ASSETS[t.troop_type] || BUILDING_ASSETS.DEFAULT;
                        const isMarching = battlePhase === 'simulating';
                        return (
                            <img key={t.id} src={src} alt={t.troop_type}
                                style={{
                                    position: 'absolute',
                                    top: isMarching ? `${t.targetYPct || 50}%` : `${t.yPct}%`,
                                    left: isMarching ? `${t.targetXPct || 50}%` : `${t.xPct}%`,
                                    transform: 'translate(-50%, -50%)',
                                    width: '38px', height: '38px',
                                    objectFit: 'contain',
                                    zIndex: 700,
                                    pointerEvents: 'none',
                                    filter: 'drop-shadow(0 3px 6px rgba(0,0,0,0.9))',
                                    transition: isMarching ? 'top 1.8s ease-in, left 1.8s ease-in, opacity 1.8s ease-in, width 1.8s, height 1.8s' : 'none',
                                    opacity: isMarching ? 0 : 1
                                }}
                            />
                        );
                    })}
                </div>
            </div>


            {battlePhase !== 'finished' && (
                <div style={{
                    background: 'linear-gradient(to top, #1a0e05cc, #3a2210cc)',
                    backdropFilter: 'blur(6px)',
                    borderTop: '3px solid #1a1007',
                    padding: '14px 20px',
                    display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                    flexShrink: 0, minHeight: '110px'
                }}>
                    <div style={{ display: 'flex', gap: '10px', overflowX: 'auto', flex: 1, marginRight: '16px', alignItems: 'center' }}>
                        {army.length === 0 ? (
                            <span style={{ color: '#aaa' }}>No trained troops. Train some first!</span>
                        ) : army.map((t, idx) => {
                            const key = `${t.troop_type}_${t.current_level}`;
                            const isSelected = selectedTroopKey === key;
                            const src = TROOP_ASSETS[t.troop_type] || BUILDING_ASSETS.DEFAULT;
                            const remaining = deployments[key] || 0;

                            return (
                                <div
                                    key={idx}
                                    onClick={() => battlePhase === 'deploying' && setSelectedTroopKey(key)}
                                    style={{
                                        background: isSelected ? 'rgba(80, 200, 80, 0.3)' : 'rgba(0,0,0,0.55)',
                                        border: isSelected ? '2px solid #5cdb5c' : '2px solid #a67c52',
                                        borderRadius: '8px',
                                        padding: '8px 10px', color: '#fff', textAlign: 'center', minWidth: '80px', flexShrink: 0,
                                        cursor: battlePhase === 'deploying' ? 'pointer' : 'default',
                                        opacity: remaining <= 0 ? 0.4 : 1,
                                        transition: 'all 0.2s'
                                    }}
                                >
                                    <img src={src} alt={t.troop_type} style={{ width: '36px', height: '36px', objectFit: 'contain' }} />
                                    <div style={{ fontSize: '11px', fontWeight: 'bold', textTransform: 'capitalize', marginTop: '3px' }}>
                                        {t.troop_type} Lv{t.current_level}
                                    </div>
                                    <div style={{ color: '#ffce00', fontSize: '14px', fontWeight: 'bold', marginTop: '4px' }}>
                                        x{remaining}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                    <button
                        onClick={handleAttack}
                        disabled={battlePhase !== 'deploying' || gridDeployedTroops.length === 0}
                        style={{
                            padding: '14px 26px', fontSize: '20px', fontWeight: '900',
                            background: 'linear-gradient(to bottom, #ff5e5e, #cc0000)',
                            border: '3px solid #660000', borderRadius: '8px', color: '#fff',
                            cursor: battlePhase !== 'deploying' ? 'not-allowed' : 'pointer',
                            boxShadow: 'inset 0 2px 0 rgba(255,255,255,0.4), 0 6px 0 #4d0000',
                            textShadow: '2px 2px 0 #4d0000',
                            opacity: battlePhase !== 'deploying' ? 0.55 : 1,
                            flexShrink: 0
                        }}
                    >
                        {battlePhase === 'simulating' ? '⚔️ ATTACKING...' : '🗡️ ATTACK!'}
                    </button>
                </div>
            )}


            {battlePhase === 'finished' && battleResult && (
                <div style={{
                    position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh',
                    backgroundColor: 'rgba(0,0,0,0.88)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 800
                }}>
                    <div style={{
                        background: 'linear-gradient(to bottom,#3a2512,#24160a)', border: '5px solid #ffce00',
                        borderRadius: '16px', width: '420px', padding: '30px', color: '#fff', textAlign: 'center',
                        boxShadow: '0 0 40px rgba(255,206,0,0.4)',
                        animation: 'dropIn 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275)'
                    }}>
                        <div style={{ fontSize: '60px', marginBottom: '10px' }}>
                            {battleResult.destruction_percent >= 100 ? '⭐⭐⭐' : battleResult.destruction_percent >= 50 ? '⭐⭐' : battleResult.destruction_percent > 0 ? '⭐' : '❌'}
                        </div>
                        <h1 style={{ fontSize: '40px', color: battleResult.victory ? '#aee536' : '#ff5e5e', textShadow: '2px 2px 0 #000', margin: '0 0 16px' }}>
                            {battleResult.victory ? 'VICTORY!' : 'DEFEAT'}
                        </h1>
                        <div style={{ fontSize: '22px', fontWeight: 'bold', marginBottom: '20px' }}>
                            Destruction: <span style={{ color: '#ffce00' }}>{battleResult.destruction_percent}%</span>
                        </div>
                        <div style={{ background: 'rgba(0,0,0,0.5)', padding: '14px', borderRadius: '8px', marginBottom: '20px', textAlign: 'left' }}>
                            <h3 style={{ margin: '0 0 10px', color: '#ffce00', textAlign: 'center' }}>Loot Acquired</h3>
                            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '17px', fontWeight: 'bold', marginBottom: '5px' }}>
                                <span>Gold:</span><span>🪙 {battleResult.loot_stolen?.gold_coins || 0}</span>
                            </div>
                            <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '17px', fontWeight: 'bold' }}>
                                <span>Elixir:</span><span>⚗️ {battleResult.loot_stolen?.elixir || 0}</span>
                            </div>
                        </div>
                        <div style={{ fontSize: '20px', fontWeight: 'bold', marginBottom: '24px' }}>
                            XP: <span style={{ color: battleResult.victory ? '#aee536' : '#ff5e5e' }}>
                                {battleResult.victory ? '+' : '-'}{battleResult.xp_modifier}
                            </span>
                        </div>
                        <button className="coc-button" style={{ margin: 0 }} onClick={onClose}>RETURN HOME</button>
                    </div>
                </div>
            )}

            <style>{`
                @keyframes shimmer {
                    0%   { background-position: 0% 50%; }
                    100% { background-position: 300% 50%; }
                }
                @keyframes bannerPulse {
                    0%   { box-shadow: 0 4px 20px rgba(255,50,0,0.5); }
                    100% { box-shadow: 0 4px 35px rgba(255,180,0,0.9); }
                }
                @keyframes shake {
                    0% { transform: translate(1px, 1px) rotate(0deg); }
                    10% { transform: translate(-1px, -2px) rotate(-1deg); }
                    20% { transform: translate(-3px, 0px) rotate(1deg); }
                    30% { transform: translate(3px, 2px) rotate(0deg); }
                    40% { transform: translate(1px, -1px) rotate(1deg); }
                    50% { transform: translate(-1px, 2px) rotate(-1deg); }
                    60% { transform: translate(-3px, 1px) rotate(0deg); }
                    70% { transform: translate(3px, 1px) rotate(-1deg); }
                    80% { transform: translate(-1px, -1px) rotate(1deg); }
                    90% { transform: translate(1px, 2px) rotate(0deg); }
                    100% { transform: translate(1px, -2px) rotate(-1deg); }
                }
                @keyframes dropIn {
                    0% { transform: translateY(-100vh) scale(0.5); opacity: 0; }
                    100% { transform: translateY(0) scale(1); opacity: 1; }
                }
            `}</style>
        </div>
    );
}
