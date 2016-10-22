package testutils

import (
	"log"
	"time"

	"github.com/codegp/cloud-persister"
	"github.com/codegp/cloud-persister/models"
	"github.com/codegp/game-object-types/types"
	"github.com/codegp/kube-client"
)

type TestUtils struct {
	cp      *cloudpersister.CloudPersister
	kclient *kubeclient.KubeClient
}

func NewTestUtils(cp *cloudpersister.CloudPersister, kclient *kubeclient.KubeClient) *TestUtils {
	return &TestUtils{
		cp:      cp,
		kclient: kclient,
	}
}

// Builds a test gametype
// Waits for the build to compelete if watch is true
// Does not rebuild if a test game type already exists unless force is true
func (u *TestUtils) BuildTestGametype(watch bool, force bool) error {
	gameType, exists, err := u.getOrCreateTestGameType()
	if err != nil {
		return err
	}

	if exists && !force {
		return nil
	}

	pod, err := u.kclient.BuildGameType(gameType)
	if err != nil {
		return err
	}
	if watch {
		_, err = u.kclient.WatchToCompletion(pod)
		if err != nil {
			return err
		}
	}

	_, err = u.getOrCreateTestMap(gameType)
	if err != nil {
		return err
	}
	
	_, err = u.getOrCreateTestProject(gameType)
	return err
}

// Runs a test game using the test game type
// Waits for the game to complete if watch is true
func (u *TestUtils) RunTestGame(watch bool) error {
	gameType, exists, err := u.getOrCreateTestGameType()
	if err != nil {
		return err
	}

	if !exists {
		err = u.BuildTestGametype(true, true)
		if err != nil {
			return err
		}
	}

	m, err := u.getOrCreateTestMap(gameType)
	if err != nil {
		return err
	}
	proj, err := u.getOrCreateTestProject(gameType)
	if err != nil {
		return err
	}
	game, err := u.createTestGame(gameType, proj, m)
	if err != nil {
		return err
	}
	pod, err := u.kclient.StartGame(game)
	if err != nil {
		return err
	}
	if watch {
		_, err = u.kclient.WatchToCompletion(pod)
	}
	return err
}

// already exists
func (u *TestUtils) getOrCreateTestGameType() (*models.GameType, bool, error) {
	gameType, err := u.cp.QueryGameTypesByProp("Name", "testGameType")
	if err != nil {
		return nil, false, err
	}

	if gameType != nil {
		log.Println("Test game type found!")
		return gameType, true, nil
	}

	log.Println("Test game type not found. Creating...")

	bot, err := u.createBotType()
	if err != nil {
		return nil, false, err
	}

	terrain, err := u.createTerrainType()
	if err != nil {
		return nil, false, err
	}

	item, err := u.createItemType()
	if err != nil {
		return nil, false, err
	}

	gameType = &models.GameType{
		Name:         "testGameType",
		Version:      "testGameType",
		BotTypes:     []int64{bot.ID},
		TerrainTypes: []int64{terrain.ID},
		ItemTypes:    []int64{item.ID},
		ApiFuncs:     []string{"me", "canMove", "move"},
		NumTeams:     1,
		CreatorID:    12345,
		Description:  "testGameType",
	}

	gameType, err = u.cp.AddGameType(gameType)
	if err != nil {
		return nil, false, err
	}
	testCode, err := Asset("testfiles/testgametype.go.tmpl")
	if err != nil {
		return nil, false, err
	}
	err = u.cp.WriteGameTypeCode(gameType.ID, testCode)
	if err != nil {
		return nil, false, err
	}
	return gameType, false, nil
}

func (u *TestUtils) createBotType() (*types.BotType, error) {
	log.Println("Creating move type...")
	moveType := &types.MoveType{
		Name:  "testMoveType",
		Delay: 0,
		TakesDelayFromTerrain: false,
	}
	moveType, err := u.cp.AddMoveType(moveType)
	if err != nil {
		return nil, err
	}

	attackType := &types.AttackType{
		Name:             "testAttackType",
		Damage:           1,
		Delay:            0,
		Range:            0,
		Accuracy:         1,
		AttackDelayDealt: 0,
		MoveDelayDealt:   0,
	}

	attackType, err = u.cp.AddAttackType(attackType)
	if err != nil {
		return nil, err
	}

	log.Println("Creating bot type...")
	botType := &types.BotType{
		Name:              "testBot",
		AttackTypeIDs:     []int64{attackType.ID},
		MoveTypeIDs:       []int64{moveType.ID},
		CanSpawn:          true,
		CanBeSpawned:      true,
		SpawnDelay:        1,
		MaxHealth:         100,
		CanHeal:           true,
		MoveDelayFactor:   1,
		DamageFactor:      1,
		AttackDelayFactor: 1,
		RangeFactor:       1,
		AccuracyFactor:    1,
		SpawnDelayFactor:  1,
	}

	botType, err = u.cp.AddBotType(botType)
	if err != nil {
		return nil, err
	}
	icon, err := Asset("testfiles/bot.png")
	if err != nil {
		return nil, err
	}

	err = u.cp.WriteIcon(botType.ID, icon)
	if err != nil {
		return nil, err
	}

	return botType, nil
}

func (u *TestUtils) createTerrainType() (*types.TerrainType, error) {
	log.Println("Creating terrain type...")
	terrainType := &types.TerrainType{
		Name:            "testTerrain",
		CanBeOccupied:   true,
		MoveDelayFactor: 1,
		DamagePenalty:   1,
	}

	terrainType, err := u.cp.AddTerrainType(terrainType)
	if err != nil {
		return nil, err
	}
	icon, err := Asset("testfiles/dirt.png")
	if err != nil {
		return nil, err
	}

	err = u.cp.WriteIcon(terrainType.ID, icon)
	if err != nil {
		return nil, err
	}

	return terrainType, nil
}

func (u *TestUtils) createItemType() (*types.ItemType, error) {
	log.Println("Creating item type...")
	itemType := &types.ItemType{
		Name:              "testItem",
		MoveDelayFactor:   1,
		DamageFactor:      1,
		AttackDelayFactor: 1,
		RangeFactor:       1,
		AccuracyFactor:    1,
		SpawnDelayFactor:  1,
	}

	itemType, err := u.cp.AddItemType(itemType)
	if err != nil {
		return nil, err
	}
	icon, err := Asset("testfiles/flag.png")
	if err != nil {
		return nil, err
	}

	err = u.cp.WriteIcon(itemType.ID, icon)
	if err != nil {
		return nil, err
	}

	return itemType, nil
}

func (u *TestUtils) getOrCreateTestMap(gameType *models.GameType) (*models.Map, error) {
	m, err := u.cp.QueryMapsByProp("Name", "testMap")
	if err != nil {
		return nil, err
	}
	if m != nil {
		log.Println("Test map found!")
		return m, nil
	}

	log.Println("Test map not found. Creating...")

	m = &models.Map{
		Name:       "testMap",
		GameTypeID: gameType.ID,
		RoundLimit: 10,
	}

	m, err = u.cp.AddMap(m)
	if err != nil {
		return nil, err
	}

	testMap, err := Asset("testfiles/testmap.json")
	if err != nil {
		return nil, err
	}

	err = u.cp.WriteMap(m.ID, testMap)
	if err != nil {
		return nil, err
	}

	gameType.MapIDs = append(gameType.MapIDs, m.ID)
	err = u.cp.UpdateGameType(gameType)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (u *TestUtils) getOrCreateTestProject(gameType *models.GameType) (*models.Project, error) {
	proj, err := u.cp.QueryProjectsByProp("Name", "testProject")
	if err != nil {
		return nil, err
	}
	if proj != nil {
		log.Println("Test project found!")
		return proj, nil
	}

	log.Println("Test project not found. Creating...")

	proj = &models.Project{
		Name:       "testProject",
		Language:   "go",
		GameTypeID: gameType.ID,
		UserID:     12345,
		FileNames:  []string{"testbot"},
		GameIDs:    []int64{},
	}

	proj, err = u.cp.AddProject(proj)
	if err != nil {
		return nil, err
	}

	testCode, err := Asset("testfiles/testbot.go.tmpl")
	if err != nil {
		return nil, err
	}

	err = u.cp.WriteProjectFile(proj.ID, "testbot", testCode)
	if err != nil {
		return nil, err
	}

	user, err := u.cp.GetUser(12345)
	if err != nil {
		err = u.cp.UpdateUser(&models.User{
			ID:         12345,
			ProjectIDs: []int64{proj.ID},
		})

		if err != nil {
			return nil, err
		}

		return proj, nil
	}

	user.ProjectIDs = append(user.ProjectIDs, proj.ID)
	err = u.cp.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

func (u *TestUtils) createTestGame(gameType *models.GameType, proj *models.Project, m *models.Map) (*models.Game, error) {
	log.Printf("Starting test game...")
	game := &models.Game{
		Created:    time.Now(),
		GameTypeID: gameType.ID,
		MapID:      m.ID,
		ProjectIDs: []int64{proj.ID},
		Complete:   false,
	}

	game, err := u.cp.AddGame(game)
	if err != nil {
		return nil, err
	}
	proj.GameIDs = append(proj.GameIDs, game.ID)
	err = u.cp.UpdateProject(proj)
	if err != nil {
		return nil, err
	}

	return game, nil
}
