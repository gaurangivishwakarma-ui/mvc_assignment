import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getVillageProfile, getShopCatalog, getArmyCatalog, getArmyStatus, trainTroops, purchaseBuilding, upgradeVillage, getVillageUpgradeCost, upgradeBuilding, getBuildingUpgradeCost, completeBuildingUpgrade, moveBuilding, getMatch, collectResources } from '../api/village';

import archerImg from '../assets/troops/archer.png';
import barbarianImg from '../assets/troops/barbarian.png';
import giantImg from '../assets/troops/giant.png';
import goblinImg from '../assets/troops/goblin.png';
import wizardImg from '../assets/troops/wizard.png';

const TROOP_ASSETS = {
    archer: archerImg,
    barbarian: barbarianImg,
    giant: giantImg,
    goblin: goblinImg,
    wizard: wizardImg
};

import villageBg from '../assets/villagebg.png';

// Town Hall
import townHallShopImg from '../assets/shop/townhall.webp';
import townHall2Img from '../assets/shop/Town_Hall2.webp';
import townHall3Img from '../assets/shop/Town_Hall3.webp';
import townHall4Img from '../assets/shop/Town_Hall4.webp';

// Gold Mine
import goldMineShopImg from '../assets/shop/Gold_Mine.webp';
import goldMine2ShopImg from '../assets/shop/Gold_Mine2.webp';
import goldMine3Img from '../assets/shop/Gold_Mine3.webp';
import goldMine4Img from '../assets/shop/Gold_Mine4.webp';

// Elixir Collector
import elixirCollectorShopImg from '../assets/shop/Elixir_Collector.webp';
import elixirCollector2Img from '../assets/shop/Elixir_Collector2.webp';
import elixirCollector3Img from '../assets/shop/Elixir_Collector3.webp';
import elixirCollector4Img from '../assets/shop/Elixir_Collector4.webp';

// Gold Storage
import goldStorageShopImg from '../assets/shop/Gold_Storage.webp';
import goldStorage2Img from '../assets/shop/Gold_Storage2.webp';
import goldStorage3Img from '../assets/shop/Gold_Storage3.webp';
import goldStorage4Img from '../assets/shop/Gold_Storage4.webp';

// Elixir Storage
import elixirStorageShopImg from '../assets/shop/Elixir_Storage.webp';
import elixirStorage2Img from '../assets/shop/Elixir_Storage2.webp';
import elixirStorage3Img from '../assets/shop/Elixir_Storage3.webp';
import elixirStorage4Img from '../assets/shop/Elixir_Storage4.webp';

// Army Camp
import armyCampShopImg from '../assets/shop/Army_Camp.webp';
import armyCamp2ShopImg from '../assets/shop/Army_Camp_2.webp';
import armyCamp3Img from '../assets/shop/Army_Camp3.webp';
import armyCamp4Img from '../assets/shop/Army_Camp4.webp';

// Cannon
import canonShopImg from '../assets/shop/canon.webp';
import cannon2Img from '../assets/shop/Cannon2.webp';
import cannon3Img from '../assets/shop/Cannon3.webp';
import cannon4Img from '../assets/shop/Cannon4.webp';

// Archer Tower
import archerTowerShopImg from '../assets/shop/archer_tower.webp';
import archerTower2Img from '../assets/shop/Archer_Tower2.webp';
import archerTower3Img from '../assets/shop/Archer_Tower3.webp';
import archerTower4Img from '../assets/shop/Archer_Tower4.webp';

// Mortar
import mortarShopImg from '../assets/shop/mortar.webp';
import mortar2Img from '../assets/shop/Mortar2.webp';
import mortar3Img from '../assets/shop/Mortar3.webp';
import mortar4Img from '../assets/shop/Mortar4.webp';

import BattleOverlay from '../components/BattleOverlay';

const SHOP_ASSETS = {
    town_hall: townHallShopImg,
    town_hall_2: townHall2Img,
    town_hall_3: townHall3Img,
    town_hall_4: townHall4Img,

    gold_mine: goldMineShopImg,
    gold_mine_2: goldMine2ShopImg,
    gold_mine_3: goldMine3Img,
    gold_mine_4: goldMine4Img,

    elixir_collector: elixirCollectorShopImg,
    elixir_collector_2: elixirCollector2Img,
    elixir_collector_3: elixirCollector3Img,
    elixir_collector_4: elixirCollector4Img,

    gold_storage: goldStorageShopImg,
    gold_storage_2: goldStorage2Img,
    gold_storage_3: goldStorage3Img,
    gold_storage_4: goldStorage4Img,

    elixir_storage: elixirStorageShopImg,
    elixir_storage_2: elixirStorage2Img,
    elixir_storage_3: elixirStorage3Img,
    elixir_storage_4: elixirStorage4Img,

    army_camp: armyCampShopImg,
    army_camp_2: armyCamp2ShopImg,
    army_camp_3: armyCamp3Img,
    army_camp_4: armyCamp4Img,

    cannon: canonShopImg,
    cannon_2: cannon2Img,
    cannon_3: cannon3Img,
    cannon_4: cannon4Img,

    archer_tower: archerTowerShopImg,
    archer_tower_2: archerTower2Img,
    archer_tower_3: archerTower3Img,
    archer_tower_4: archerTower4Img,

    mortar: mortarShopImg,
    mortar_2: mortar2Img,
    mortar_3: mortar3Img,
    mortar_4: mortar4Img,
};

const BUILDING_ASSETS = {
    TOWN_HALL: '/assets/buildings/town_hall.png',
    GOLD_MINE: '/assets/buildings/gold_mine.png',
    ELIXIR_COLLECTOR: '/assets/buildings/elixir_pump.png',
    BARRACKS: '/assets/buildings/barracks.png',
    DEFAULT: '/assets/buildings/construction_placeholder.png'
};

const styles = {
    container: {
        backgroundImage: `url(${villageBg})`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
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
        background: 'linear-gradient(rgba(45, 25, 10, 0.65), rgba(45, 25, 10, 0.65)), url(/assets/sidebar_bg.png?v=5) center/cover',
        borderRight: '4px solid #1a120b',
        boxShadow: '4px 0px 10px rgba(0,0,0,0.5)',
        zIndex: 10
    },
    playerInfo: { width: '100%', marginBottom: '15px' },
    plaqueContainer: {
        background: 'linear-gradient(to bottom, #6b635e, #4a4440)',
        border: '3px solid #332d29',
        borderRadius: '6px',
        boxShadow: 'inset 0px 2px 0px rgba(255,255,255,0.2), inset 0px -2px 0px rgba(0,0,0,0.4), 0px 6px 10px rgba(0,0,0,0.7), 0px 0px 0px 2px #221d19',
        padding: '10px 20px',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        width: '100%',
        boxSizing: 'border-box'
    },
    plaqueTitle: {
        color: '#fff',
        fontWeight: '900',
        fontSize: '18px',
        fontFamily: "'Montserrat', sans-serif",
        textShadow: '1px 1px 2px #000',
        letterSpacing: '1px',
        borderBottom: '2px solid #332d29',
        width: '100%',
        textAlign: 'center',
        paddingBottom: '8px',
        marginBottom: '8px',
        boxShadow: '0px 2px 0px rgba(255,255,255,0.1)',
        textTransform: 'uppercase'
    },
    plaqueSubtitle: {
        color: '#f0e6d2',
        fontWeight: '800',
        fontSize: '14px',
        fontFamily: "'Montserrat', sans-serif",
        textShadow: '1px 1px 2px #000'
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
    logoutBtn: { width: '150px', padding: '15px', fontSize: '20px', margin: '0' },
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
        boxSizing: 'border-box',
        cursor: 'pointer'
    },

    mapViewport: {
        flex: 1,
        overflow: 'hidden',
        padding: '20px',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center'
    },
    gridBoard: {
        display: 'grid',

        gridTemplateColumns: 'repeat(25, 1fr)',
        gridTemplateRows: 'repeat(25, 1fr)',
        gap: '1px',


        width: '90vh',
        maxWidth: '100%',
        aspectRatio: '1 / 1',

        borderRadius: '8px',
        border: '6px solid #2e4712',
        boxShadow: '0px 12px 30px rgba(0,0,0,0.4)'
    },
    buildingWrapper: {
        position: 'relative',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        cursor: 'pointer',
        transition: 'transform 0.1s ease',
        zIndex: 5
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

    bottomBar: { display: 'flex', flexDirection: 'column', justifyContent: 'space-between', gap: '15px', padding: '20px', background: 'rgba(0,0,0,0.4)', backdropFilter: 'blur(5px)', borderLeft: '2px solid rgba(255,255,255,0.1)', zIndex: 10 },
    actionButton: { width: '100%', padding: '15px', fontSize: '18px', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px', margin: 0 },
    attackButton: { background: 'linear-gradient(to bottom, #ff5e5e, #cc0000)', borderColor: '#660000', boxShadow: 'inset 0px 2px 0px rgba(255,255,255,0.4), 0px 6px 0px #4d0000', textShadow: '2px 2px 0px #4d0000' },
    loadingScreen: { flex: 1, display: 'flex', justifyContent: 'center', alignItems: 'center', background: '#315800', color: '#fff', fontSize: '24px', fontWeight: '800' },

    modalOverlay: { position: 'absolute', top: 0, left: 0, right: 0, bottom: 0, backgroundColor: 'rgba(0,0,0,0.75)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 100, padding: '20px' },
    modalBox: { background: 'linear-gradient(to bottom, #3a2512, #24160a)', border: '5px solid #1a1007', borderRadius: '16px', width: '90%', maxWidth: '500px', padding: '20px', boxShadow: '0px 10px 25px rgba(0,0,0,0.8)', color: '#fff' },
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
    const [hoveredBuilding, setHoveredBuilding] = useState(null);
    const [isArmyOpen, setIsArmyOpen] = useState(false);
    const [armyCatalog, setArmyCatalog] = useState([]);
    const [armyStatus, setArmyStatus] = useState([]);
    const [hoveredTroop, setHoveredTroop] = useState(null);
    const [troopQuantities, setTroopQuantities] = useState({});
    const [isDashboardModalOpen, setIsDashboardModalOpen] = useState(false);
    const [activeDashboardTab, setActiveDashboardTab] = useState('resources');
    const [selectedBuildingToPlace, setSelectedBuildingToPlace] = useState(null);
    const [selectedBuildingToMove, setSelectedBuildingToMove] = useState(null);
    const [isMoveModalOpen, setIsMoveModalOpen] = useState(false);
    const [isVillageUpgradeModalOpen, setIsVillageUpgradeModalOpen] = useState(false);
    const [villageUpgradeCostData, setVillageUpgradeCostData] = useState(null);
    const [upgradingTimers, setUpgradingTimers] = useState({});
    const [buildingUpgradeModal, setBuildingUpgradeModal] = useState(null);
    const [matchData, setMatchData] = useState(null);
    const [isMatchModalOpen, setIsMatchModalOpen] = useState(false);
    const [isBattleActive, setIsBattleActive] = useState(false);
    const [collectionFloating, setCollectionFloating] = useState(null);

    useEffect(() => {
        const interval = setInterval(() => {
            setUpgradingTimers(prev => {
                const next = { ...prev };
                let hasChanges = false;
                for (const id in next) {
                    if (next[id] > 0) {
                        next[id] -= 1;
                        hasChanges = true;

                        if (next[id] === 0) {
                            completeBuildingUpgrade({ placement_id: id }).then(() => {
                                getVillageProfile().then(setPlayer);
                            }).catch(err => alert(`Completion failed: ${err}`));
                        }
                    }
                }
                return hasChanges ? next : prev;
            });
        }, 1000);
        return () => clearInterval(interval);
    }, []);

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
        setSelectedBuildingToPlace(item);
        setIsShopOpen(false);
        alert(`Select a tile on the grid to place your ${item.name}!`);
    };

    const handleTileClick = async (x, y) => {
        if (selectedBuildingToMove) {
            try {
                const result = await moveBuilding({
                    owned_building_id: selectedBuildingToMove.placement_id,
                    new_x: x,
                    new_y: y
                });
                alert(result.status || `Successfully moved building!`);
                const profileData = await getVillageProfile();
                setPlayer(profileData);
            } catch (err) {
                alert(`Move Failed: ${err}`);
            } finally {
                setSelectedBuildingToMove(null);
            }
            return;
        }

        if (!selectedBuildingToPlace) return;

        try {
            const result = await purchaseBuilding({
                building_id: selectedBuildingToPlace.id,
                x_coords: x,
                y_coords: y
            });
            alert(result.message || `Successfully placed ${selectedBuildingToPlace.name}!`);

            if (result.building && result.building.id) {
                setUpgradingTimers(prev => ({ ...prev, [result.building.id]: 5 }));
            }

            const profileData = await getVillageProfile();
            setPlayer(profileData);
        } catch (err) {
            alert(`Placement Failed: ${err}`);
        } finally {
            setSelectedBuildingToPlace(null);
        }
    };

    const handleQtyChange = (troopType, troopLevel, delta) => {
        const key = `${troopType}_${troopLevel}`;
        setTroopQuantities(prev => {
            const current = prev[key] || 1;
            return { ...prev, [key]: Math.max(1, current + delta) };
        });
    };

    const handleTrainTroop = async (troop) => {
        const key = `${troop.troop_type}_${troop.level_req}`;
        const qty = troopQuantities[key] || 1;
        try {
            const result = await trainTroops({
                troop_type: troop.troop_type,
                level: troop.level_req,
                quantity: qty
            });
            alert(result.message);

            const [profileData, statusData] = await Promise.all([
                getVillageProfile(),
                getArmyStatus()
            ]);
            setPlayer(profileData);
            setArmyStatus(statusData.army || []);

            setTroopQuantities(prev => ({ ...prev, [key]: 1 }));
            setIsArmyOpen(false);
        } catch (err) {
            alert(`Training Failed: ${err}`);
        }
    };

    const handleUpgradeVillageClick = async () => {
        try {
            const costData = await getVillageUpgradeCost();
            setVillageUpgradeCostData(costData);
            setIsVillageUpgradeModalOpen(true);
        } catch (err) {
            alert(`Failed to fetch upgrade info: ${err}`);
        }
    };

    const handleConfirmUpgradeVillage = async () => {
        try {
            const result = await upgradeVillage();
            alert(result.message || "Town Hall upgraded successfully!");
            const data = await getVillageProfile();
            setPlayer(data);
            setIsVillageUpgradeModalOpen(false);
        } catch (err) {
            alert(`Upgrade Failed: ${err}`);
        }
    };

    const handleUpgradeBuildingClick = async (building) => {
        try {
            const costData = await getBuildingUpgradeCost(building.placement_id);
            if (costData.is_max_level) {
                alert(costData.message || "Building is already at max level!");
                return;
            }
            setBuildingUpgradeModal({
                building,
                nextLevel: costData.next_level,
                cost: costData.cost,
                costType: costData.cost_type,
                name: costData.name || building.building_type?.replace(/_/g, ' ')
            });
        } catch (err) {
            alert(`Failed to fetch upgrade cost: ${err}`);
        }
    };

    const handleConfirmBuildingUpgrade = async () => {
        const { building } = buildingUpgradeModal;
        setBuildingUpgradeModal(null);
        try {
            const result = await upgradeBuilding({ placement_id: building.placement_id });
            alert(result.message || 'Upgrade started!');
            setUpgradingTimers(prev => ({ ...prev, [building.placement_id]: 5 }));
            const data = await getVillageProfile();
            setPlayer(data);
        } catch (err) {
            alert(`Upgrade Failed: ${err}`);
        }
    };

    const handleMatchClick = async () => {
        try {
            const data = await getMatch();
            setMatchData(data);
            setIsMatchModalOpen(true);
        } catch (err) {
            alert(`Matchmaking Failed: ${err}`);
        }
    };

    const handleCollectResources = async () => {
        try {
            const data = await collectResources();
            if (data.gold_looted > 0 || data.elixir_looted > 0) {
                setCollectionFloating({ gold: data.gold_looted, elixir: data.elixir_looted });
                setTimeout(() => setCollectionFloating(null), 2000);
            } else {
                alert(data.message);
            }
            const profileData = await getVillageProfile();
            setPlayer(profileData);
        } catch (err) {
            alert(`Collection Failed: ${err}`);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/login');
    };

    if (error) return <div style={styles.container}> <div style={styles.loadingScreen}>⚔️ {error}</div> </div>;
    if (!player) return <div style={styles.container}> <div style={styles.loadingScreen}>LOADING BASE...</div> </div>;

    const deployedBuildings = player.buildings_placed || [];

    const townHallObj = {
        building_type: 'town_hall',
        name: `Town Hall Lv${player.village_level || 1}`,
        x_coords: 11,
        y_coords: 11,
        width: 4,
        breadth: 4,
        current_level: player.village_level || 1
    };

    const allBuildings = [townHallObj, ...deployedBuildings];

    return (
        <div style={styles.container}>
            <div style={styles.topBar}>
                <div style={styles.playerInfo}>
                    <div style={styles.plaqueContainer}>
                        <div style={styles.plaqueTitle}>
                            <span>{player.username || 'ADMIN'}</span>
                        </div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', justifyContent: 'center', marginBottom: '4px' }}>
                            <span style={{ background: 'linear-gradient(to bottom, #ffe066, #c8860a)', color: '#3a1e00', fontWeight: '900', fontSize: '12px', padding: '2px 8px', borderRadius: '12px', border: '2px solid #7c4e00', boxShadow: '0 2px 4px rgba(0,0,0,0.5)', letterSpacing: '0.5px' }}>
                                {player.xp_points || 0} XP
                            </span>
                        </div>
                        <div
                            style={{ ...styles.plaqueSubtitle, cursor: 'pointer', color: '#ffce00' }}
                            onClick={handleUpgradeVillageClick}
                            title="Click to upgrade Town Hall"
                        >
                            Town Hall Level {player.village_level || 1}
                        </div>
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', gap: '15px', marginTop: '15px', width: '100%', boxSizing: 'border-box' }}>
                        <button
                            className="coc-button"
                            style={styles.actionButton}
                            onClick={() => setIsDashboardModalOpen(true)}
                        >
                            DASHBOARD
                        </button>

                        <button className="coc-button" style={styles.actionButton} onClick={handleCollectResources}>
                            COLLECT
                        </button>

                        <button className="coc-button" style={{ ...styles.actionButton, background: 'linear-gradient(to bottom, #f29f05, #d96d00)', borderColor: '#a64a00', boxShadow: 'inset 0px 2px 0px rgba(255,255,255,0.4), 0px 6px 0px #733000', textShadow: '2px 2px 0px #733000' }} onClick={handleMatchClick}>
                            MATCH
                        </button>
                    </div>
                </div>

                <div style={styles.resources}>
                    <button className="coc-button" style={styles.actionButton} onClick={handleLogout}>
                        LOGOUT
                    </button>
                </div>
            </div>


            <div style={styles.mapViewport}>
                <div style={styles.gridBoard}>
                    {Array.from({ length: 625 }).map((_, i) => {
                        const x = (i % 25) + 1;
                        const y = Math.floor(i / 25) + 1;
                        return (
                            <div
                                key={`tile-${i}`}
                                style={{
                                    gridColumn: x,
                                    gridRow: y,
                                    backgroundImage: 'url(/assets/grass_tile.png)',
                                    backgroundSize: 'cover',
                                    cursor: (selectedBuildingToPlace || selectedBuildingToMove) ? 'crosshair' : 'default',
                                    zIndex: 1
                                }}
                                onClick={() => handleTileClick(x, y)}
                            />
                        )
                    })}

                    {allBuildings.map((b, index) => {
                        const levelSuffix = b.current_level > 1 ? `_${b.current_level}` : '';
                        const typeWithLevel = `${b.building_type}${levelSuffix}`;
                        const imageSrc = SHOP_ASSETS[typeWithLevel] || SHOP_ASSETS[b.building_type] || BUILDING_ASSETS[b.building_type?.toUpperCase()] || BUILDING_ASSETS.DEFAULT;

                        const isGoldMine = b.building_type === 'gold_mine' || b.building_type === 'gold_mine_2';
                        const isElixirCollector = b.building_type === 'elixir_collector';

                        const numGoldMines = allBuildings.filter(x => x.building_type === 'gold_mine' || x.building_type === 'gold_mine_2').length || 1;
                        const numElixirCollectors = allBuildings.filter(x => x.building_type === 'elixir_collector').length || 1;

                        return (
                            <div
                                key={index}
                                style={{
                                    ...styles.buildingWrapper,
                                    gridColumn: `${b.x_coords || 1} / span ${b.width || 2}`,
                                    gridRow: `${b.y_coords || 1} / span ${b.breadth || b.height || 2}`
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

                                {collectionFloating && collectionFloating.gold > 0 && isGoldMine && (
                                    <div className="floating-resource" style={{ position: 'absolute', top: '0px', color: '#ffce00', fontWeight: '900', fontSize: '24px', textShadow: '2px 2px 4px rgba(0,0,0,0.8), -1px -1px 0 #000, 1px -1px 0 #000, -1px 1px 0 #000, 1px 1px 0 #000', pointerEvents: 'none', zIndex: 50 }}>
                                        +{Math.floor(collectionFloating.gold / numGoldMines)} 🪙
                                    </div>
                                )}
                                {collectionFloating && collectionFloating.elixir > 0 && isElixirCollector && (
                                    <div className="floating-resource" style={{ position: 'absolute', top: '0px', color: '#ff5eff', fontWeight: '900', fontSize: '24px', textShadow: '2px 2px 4px rgba(0,0,0,0.8), -1px -1px 0 #000, 1px -1px 0 #000, -1px 1px 0 #000, 1px 1px 0 #000', pointerEvents: 'none', zIndex: 50 }}>
                                        +{Math.floor(collectionFloating.elixir / numElixirCollectors)} ⚗️
                                    </div>
                                )}
                            </div>
                        );
                    })}
                </div>
            </div>

            <div style={styles.bottomBar}>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '15px', width: '100%' }}>
                    <button className="coc-button" style={styles.actionButton} onClick={() => setIsShopOpen(true)}>
                        SHOP
                    </button>
                    <button className="coc-button" style={styles.actionButton} onClick={() => setIsMoveModalOpen(true)}>
                        MOVE
                    </button>
                </div>
                <button className="coc-button" style={styles.actionButton} onClick={() => setIsArmyOpen(true)}>
                    ARMY
                </button>
            </div>

            {isMatchModalOpen && matchData && (
                <div style={styles.modalOverlay}>
                    <div style={{ ...styles.modalBox, width: '350px', textAlign: 'center' }}>
                        <div style={styles.modalHeader}>
                            <h2 style={{ margin: 0, fontWeight: '900' }}>OPPONENT FOUND</h2>
                            <button className="coc-button" style={{ width: '40px', margin: 0, padding: '5px' }} onClick={() => setIsMatchModalOpen(false)}>X</button>
                        </div>
                        <div style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '15px', alignItems: 'center' }}>
                            <div style={{ fontSize: '22px', fontWeight: 'bold', color: '#ffce00' }}>
                                {matchData.username}
                            </div>
                            <div style={{ fontSize: '16px' }}>
                                Town Hall Level: {matchData.village_level} | XP: {matchData.xp_points}
                            </div>
                            <div style={{ width: '100%', background: 'rgba(0,0,0,0.5)', padding: '15px', borderRadius: '8px', marginTop: '10px' }}>
                                <h3 style={{ margin: '0 0 10px 0', color: '#ffce00' }}>Available Loot</h3>
                                <div style={{ display: 'flex', justifyContent: 'space-around', fontSize: '18px', fontWeight: 'bold' }}>
                                    <span>🪙 {matchData.loot_available?.gold_coins || 0}</span>
                                    <span>⚗️ {matchData.loot_available?.elixir || 0}</span>
                                </div>
                            </div>
                            <div style={{ display: 'flex', gap: '10px', width: '100%', marginTop: '15px' }}>
                                <button className="coc-button" style={{ flex: 1, padding: '12px' }} onClick={() => {
                                    setIsMatchModalOpen(false);
                                    handleMatchClick(); // Search again
                                }}>
                                    NEXT
                                </button>
                                <button className="coc-button" style={{ flex: 1, padding: '12px', ...styles.attackButton }} onClick={() => {
                                    setIsMatchModalOpen(false);
                                    setIsBattleActive(true);
                                }}>
                                    ATTACK
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {isBattleActive && matchData && (
                <BattleOverlay
                    matchData={matchData}
                    onClose={async () => {
                        setIsBattleActive(false);
                        const profileData = await getVillageProfile();
                        setPlayer(profileData);
                    }}
                    BUILDING_ASSETS={BUILDING_ASSETS}
                    SHOP_ASSETS={SHOP_ASSETS}
                    TROOP_ASSETS={TROOP_ASSETS}
                />
            )}

            {isMoveModalOpen && (
                <div style={styles.modalOverlay}>
                    <div style={{ ...styles.modalBox, width: '400px' }}>
                        <div style={styles.modalHeader}>
                            <h2 style={{ margin: 0, fontWeight: '900' }}>MOVE BUILDING</h2>
                            <button className="coc-button" style={{ width: '40px', margin: 0, padding: '5px' }} onClick={() => setIsMoveModalOpen(false)}>X</button>
                        </div>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', maxHeight: '300px', overflowY: 'auto' }}>
                            {deployedBuildings.length > 0 ? deployedBuildings.map((b, i) => (
                                <div key={i} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'rgba(255,255,255,0.1)', padding: '8px', borderRadius: '4px' }}>
                                    <div>
                                        <div style={{ textTransform: 'capitalize' }}>{b.name || b.building_type?.replace(/_/g, ' ')}</div>
                                        <div style={{ color: '#aaa', fontSize: '12px' }}>Lv {b.current_level || 1} | Pos: {b.x_coords},{b.y_coords}</div>
                                    </div>
                                    {b.placement_id && (
                                        <button
                                            className="coc-button"
                                            style={{ width: 'auto', padding: '5px 10px', fontSize: '12px', margin: 0 }}
                                            onClick={() => {
                                                setSelectedBuildingToMove(b);
                                                setIsMoveModalOpen(false);
                                                alert(`Select a new tile on the grid to move your ${b.name || b.building_type?.replace(/_/g, ' ')}!`);
                                            }}
                                        >
                                            Move
                                        </button>
                                    )}
                                </div>
                            )) : <div style={{ textAlign: 'center', color: '#aaa' }}>No buildings placed yet.</div>}
                        </div>
                    </div>
                </div>
            )}

            {isVillageUpgradeModalOpen && (
                <div style={styles.modalOverlay}>
                    <div style={{ ...styles.modalBox, width: '350px', textAlign: 'center' }}>
                        <div style={styles.modalHeader}>
                            <h2 style={{ margin: 0, fontWeight: '900' }}>TOWN HALL UPGRADE</h2>
                            <button className="coc-button" style={{ width: '40px', margin: 0, padding: '5px' }} onClick={() => setIsVillageUpgradeModalOpen(false)}>X</button>
                        </div>
                        {villageUpgradeCostData?.is_max_level ? (
                            <div style={{ padding: '20px', fontSize: '18px', color: '#ffce00', fontWeight: 'bold' }}>
                                {villageUpgradeCostData.message}
                            </div>
                        ) : (
                            <div style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '15px', alignItems: 'center' }}>
                                <div style={{ fontSize: '18px' }}>
                                    Upgrade to <span style={{ fontWeight: 'bold', color: '#ffce00' }}>Level {villageUpgradeCostData?.next_level}</span>?
                                </div>
                                <div style={{ fontSize: '20px', fontWeight: 'bold' }}>
                                    Cost: <span style={{ color: '#ffce00' }}>{villageUpgradeCostData?.cost_type === 'elixir' ? '⚗️' : '🪙'} {villageUpgradeCostData?.cost}</span>
                                </div>
                                <button className="coc-button" style={{ ...styles.actionButton, marginTop: '10px' }} onClick={handleConfirmUpgradeVillage}>
                                    Confirm Upgrade
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            )}

            {isShopOpen && (
                <div style={styles.modalOverlay}>
                    <div style={{ ...styles.modalBox, position: 'relative' }}>
                        {hoveredBuilding && (
                            <div style={{
                                position: 'absolute',
                                right: '105%',
                                top: '0',
                                background: 'linear-gradient(to bottom, #4a3525, #2d1f14)',
                                border: '3px solid #7c6248',
                                borderRadius: '12px',
                                padding: '15px',
                                width: '220px',
                                color: '#fff',
                                boxShadow: '0px 10px 25px rgba(0,0,0,0.8)',
                                zIndex: 110,
                                display: 'flex',
                                flexDirection: 'column',
                                gap: '10px'
                            }}>
                                <div style={{ fontWeight: '900', fontSize: '18px', textTransform: 'capitalize', color: '#ffce00', textAlign: 'center' }}>
                                    {hoveredBuilding.name}
                                </div>
                                <div style={{ display: 'flex', justifyContent: 'center', background: '#fff', borderRadius: '8px', padding: '5px' }}>
                                    <img
                                        src={SHOP_ASSETS[hoveredBuilding.building_type] || BUILDING_ASSETS.DEFAULT}
                                        alt={hoveredBuilding.name}
                                        style={{ width: '100%', height: '180px', objectFit: 'contain', borderRadius: '4px' }}
                                        onError={(e) => {
                                            e.target.style.display = 'none';
                                        }}
                                    />
                                </div>
                                <div style={{ background: 'rgba(0,0,0,0.5)', padding: '10px', borderRadius: '8px', fontSize: '13px', display: 'flex', flexDirection: 'column', gap: '5px' }}>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>HP:</span> <span style={{ fontWeight: 'bold' }}>{hoveredBuilding.hit_points}</span></div>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>Footprint:</span> <span style={{ fontWeight: 'bold' }}>{hoveredBuilding.width}x{hoveredBuilding.breadth}</span></div>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>Required TH:</span> <span style={{ fontWeight: 'bold' }}>Lv {hoveredBuilding.level_req}</span></div>
                                </div>
                            </div>
                        )}
                        <div style={styles.modalHeader}>
                            <h2 style={{ margin: 0, fontWeight: '900' }}>VILLAGE SHOP</h2>
                            <button className="coc-button" style={{ width: '40px', margin: 0, padding: '5px' }} onClick={() => setIsShopOpen(false)}>X</button>
                        </div>
                        <div style={styles.shopGrid}>
                            {shopCatalog.map((item, idx) => (
                                <div
                                    key={idx}
                                    style={{ ...styles.shopItem, cursor: 'pointer' }}
                                    onMouseEnter={() => setHoveredBuilding(item)}
                                    onMouseLeave={() => setHoveredBuilding(null)}
                                >
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
                    <div style={{ ...styles.modalBox, position: 'relative' }}>
                        {hoveredTroop && (
                            <div style={{
                                position: 'absolute',
                                right: '105%',
                                top: '0',
                                background: 'linear-gradient(to bottom, #4a3525, #2d1f14)',
                                border: '3px solid #7c6248',
                                borderRadius: '12px',
                                padding: '15px',
                                width: '220px',
                                color: '#fff',
                                boxShadow: '0px 10px 25px rgba(0,0,0,0.8)',
                                zIndex: 110,
                                display: 'flex',
                                flexDirection: 'column',
                                gap: '10px'
                            }}>
                                <div style={{ fontWeight: '900', fontSize: '18px', textTransform: 'capitalize', color: '#ffce00', textAlign: 'center' }}>
                                    {hoveredTroop.troop_type} (Lv {hoveredTroop.level_req})
                                </div>
                                <div style={{ display: 'flex', justifyContent: 'center', background: '#fff', borderRadius: '8px', padding: '5px' }}>
                                    <img
                                        src={TROOP_ASSETS[hoveredTroop.troop_type]}
                                        alt={hoveredTroop.troop_type}
                                        style={{ width: '100%', height: '180px', objectFit: 'contain', borderRadius: '4px' }}
                                        onError={(e) => {
                                            e.target.style.display = 'none';
                                        }}
                                    />
                                </div>
                                <div style={{ fontSize: '13px', fontStyle: 'italic', color: '#ddd', textAlign: 'center' }}>
                                    "{hoveredTroop.description}"
                                </div>
                                <div style={{ background: 'rgba(0,0,0,0.5)', padding: '10px', borderRadius: '8px', fontSize: '13px', display: 'flex', flexDirection: 'column', gap: '5px' }}>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>HP:</span> <span style={{ fontWeight: 'bold' }}>{hoveredTroop.hit_points}</span></div>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>Damage:</span> <span style={{ fontWeight: 'bold' }}>{hoveredTroop.damage}</span></div>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>Speed:</span> <span style={{ fontWeight: 'bold' }}>{hoveredTroop.speed}</span></div>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>Space:</span> <span style={{ fontWeight: 'bold' }}>{hoveredTroop.housing_space}</span></div>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>Range:</span> <span style={{ fontWeight: 'bold' }}>{hoveredTroop.attack_range}</span></div>
                                </div>
                            </div>
                        )}
                        <div style={styles.modalHeader}>
                            <h2 style={{ margin: 0, fontWeight: '900' }}>TRAIN TROOPS</h2>
                            <button className="coc-button" style={{ width: '40px', margin: 0, padding: '5px' }} onClick={() => setIsArmyOpen(false)}>X</button>
                        </div>
                        <div style={styles.shopGrid}>
                            {armyCatalog.filter(troop => troop.level_req <= (player?.village_level || 1)).map((troop, idx) => (
                                <div
                                    key={idx}
                                    style={{ ...styles.shopItem, cursor: 'pointer' }}
                                    onMouseEnter={() => setHoveredTroop(troop)}
                                    onMouseLeave={() => setHoveredTroop(null)}
                                >
                                    <div>
                                        <div style={{ fontWeight: '900', fontSize: '16px', textTransform: 'capitalize' }}>{troop.troop_type} (Lv {troop.level_req})</div>
                                        <div style={{ fontSize: '12px', color: '#aaa' }}>Space: {troop.housing_space} | HP: {troop.hit_points} | Dmg: {troop.damage}</div>
                                    </div>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                        <div style={{ display: 'flex', alignItems: 'center', background: '#332211', borderRadius: '4px', overflow: 'hidden' }}>
                                            <button
                                                style={{ background: '#7c6248', color: '#fff', border: 'none', padding: '5px 10px', cursor: 'pointer', fontWeight: 'bold' }}
                                                onClick={() => handleQtyChange(troop.troop_type, troop.level_req, -1)}
                                            >-</button>
                                            <div style={{ padding: '0 5px', fontWeight: 'bold', minWidth: '30px', textAlign: 'center' }}>
                                                {troopQuantities[`${troop.troop_type}_${troop.level_req}`] || 1}
                                            </div>
                                            <button
                                                style={{ background: '#7c6248', color: '#fff', border: 'none', padding: '5px 10px', cursor: 'pointer', fontWeight: 'bold' }}
                                                onClick={() => handleQtyChange(troop.troop_type, troop.level_req, 1)}
                                            >+</button>
                                        </div>
                                        <button
                                            className="coc-button"
                                            style={{ width: 'auto', padding: '8px 15px', fontSize: '13px', margin: 0 }}
                                            onClick={() => handleTrainTroop(troop)}
                                        >
                                            Train: ⚗️ {troop.elixir_cost * (troopQuantities[`${troop.troop_type}_${troop.level_req}`] || 1)}
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            )}
            {isDashboardModalOpen && (
                <div style={styles.modalOverlay}>
                    <div style={{ ...styles.modalBox, width: '400px' }}>
                        <div style={styles.modalHeader}>
                            <h2 style={{ margin: 0, fontWeight: '900' }}>VILLAGE DASHBOARD</h2>
                            <button className="coc-button" style={{ width: '40px', margin: 0, padding: '5px' }} onClick={() => setIsDashboardModalOpen(false)}>X</button>
                        </div>
                        <div style={{ display: 'flex', gap: '10px', marginBottom: '15px' }}>
                            <button className="coc-button" style={{ flex: 1, padding: '10px', fontSize: '14px', background: activeDashboardTab === 'resources' ? '#4a8505' : '' }} onClick={() => setActiveDashboardTab('resources')}>Resources</button>
                            <button className="coc-button" style={{ flex: 1, padding: '10px', fontSize: '14px', background: activeDashboardTab === 'buildings' ? '#4a8505' : '' }} onClick={() => setActiveDashboardTab('buildings')}>Buildings</button>
                            <button className="coc-button" style={{ flex: 1, padding: '10px', fontSize: '14px', background: activeDashboardTab === 'army' ? '#4a8505' : '' }} onClick={() => setActiveDashboardTab('army')}>Army</button>
                        </div>

                        <div style={{ background: 'rgba(0,0,0,0.5)', padding: '15px', borderRadius: '8px', minHeight: '150px' }}>
                            {activeDashboardTab === 'resources' && (
                                <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', fontSize: '18px' }}>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>Gold:</span> <span style={{ fontWeight: 'bold', color: '#ffce00' }}>🪙 {player.balances?.gold || 0}</span></div>
                                    <div style={{ display: 'flex', justifyContent: 'space-between' }}><span>Elixir:</span> <span style={{ fontWeight: 'bold', color: '#ffce00' }}>⚗️ {player.balances?.elixir || 0}</span></div>
                                </div>
                            )}
                            {activeDashboardTab === 'buildings' && (
                                <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', maxHeight: '200px', overflowY: 'auto' }}>
                                    {deployedBuildings.length > 0 ? deployedBuildings.map((b, i) => (
                                        <div key={i} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'rgba(255,255,255,0.1)', padding: '8px', borderRadius: '4px' }}>
                                            <div>
                                                <div style={{ textTransform: 'capitalize' }}>{b.name || b.building_type?.replace(/_/g, ' ')}</div>
                                                <div style={{ color: '#aaa', fontSize: '12px' }}>Lv {b.current_level || 1}</div>
                                            </div>
                                            {b.placement_id && (
                                                upgradingTimers[b.placement_id] > 0 ? (
                                                    <div style={{ color: '#ffce00', fontWeight: 'bold', fontSize: '12px' }}>Upgrading... {upgradingTimers[b.placement_id]}s</div>
                                                ) : (
                                                    <button
                                                        className="coc-button"
                                                        style={{ width: 'auto', padding: '5px 10px', fontSize: '12px', margin: 0 }}
                                                        onClick={() => handleUpgradeBuildingClick(b)}
                                                    >
                                                        Upgrade ⬆
                                                    </button>
                                                )
                                            )}
                                        </div>
                                    )) : <div style={{ textAlign: 'center', color: '#aaa' }}>No buildings placed yet.</div>}

                                    {buildingUpgradeModal && (
                                        <div style={{ position: 'fixed', top: 0, left: 0, width: '100vw', height: '100vh', backgroundColor: 'rgba(0,0,0,0.85)', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: 1000 }}>
                                            <div style={{ background: 'linear-gradient(to bottom, #3a2512, #24160a)', border: '5px solid #ffce00', borderRadius: '16px', padding: '28px', width: '340px', color: '#fff', textAlign: 'center', boxShadow: '0 0 40px rgba(255,206,0,0.35)' }}>
                                                <h3 style={{ color: '#ffce00', margin: '0 0 14px', fontSize: '20px', textTransform: 'capitalize' }}>
                                                    Upgrade {buildingUpgradeModal.name}?
                                                </h3>
                                                <div style={{ fontSize: '15px', marginBottom: '8px', color: '#ccc' }}>
                                                    Level {buildingUpgradeModal.nextLevel - 1} → <span style={{ color: '#ffce00', fontWeight: 'bold' }}>Level {buildingUpgradeModal.nextLevel}</span>
                                                </div>
                                                {buildingUpgradeModal.cost !== null ? (
                                                    <div style={{ fontSize: '18px', fontWeight: 'bold', marginBottom: '20px' }}>
                                                        Cost:&nbsp;
                                                        <span style={{ color: '#ffce00' }}>
                                                            {buildingUpgradeModal.costType === 'elixir' ? '⚗️' : '🪙'}&nbsp;{buildingUpgradeModal.cost}
                                                        </span>
                                                    </div>
                                                ) : (
                                                    <div style={{ color: '#aaa', marginBottom: '20px', fontSize: '14px' }}>Max level reached or cost unavailable.</div>
                                                )}
                                                <div style={{ display: 'flex', gap: '10px', justifyContent: 'center' }}>
                                                    <button className="coc-button" style={{ margin: 0, padding: '8px 18px' }} onClick={handleConfirmBuildingUpgrade}>
                                                        ✅ Confirm
                                                    </button>
                                                    <button className="coc-button" style={{ margin: 0, padding: '8px 18px', background: 'linear-gradient(to bottom,#888,#555)', borderColor: '#333', textShadow: 'none', boxShadow: 'none' }} onClick={() => setBuildingUpgradeModal(null)}>
                                                        ✖ Cancel
                                                    </button>
                                                </div>
                                            </div>
                                        </div>
                                    )}
                                </div>
                            )}
                            {activeDashboardTab === 'army' && (
                                <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', maxHeight: '200px', overflowY: 'auto' }}>
                                    {armyStatus.length > 0 ? armyStatus.map((t, i) => (
                                        <div key={i} style={{ display: 'flex', justifyContent: 'space-between', background: 'rgba(255,255,255,0.1)', padding: '8px', borderRadius: '4px' }}>
                                            <span style={{ textTransform: 'capitalize' }}>{t.troop_type}</span>
                                            <span style={{ fontWeight: 'bold', color: '#ffce00' }}>x{t.quantity}</span>
                                        </div>
                                    )) : <div style={{ textAlign: 'center', color: '#aaa' }}>No troops trained yet.</div>}
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}