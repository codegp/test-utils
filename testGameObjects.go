package testutils

import (
  "github.com/codegp/game-object-types/types"
	"github.com/codegp/cloud-persister/models"
)

const unitTestTestID = 1

func UnitTestBotType() *types.BotType {
  botType := testBotType(unitTestTestID, unitTestTestID)
  botType.ID = unitTestTestID
  return botType
}

func UnitTestAttackType() *types.AttackType {
  attackType := testAttackType()
  attackType.ID = unitTestTestID
  return attackType
}

func UnitTestMoveType() *types.MoveType {
  moveType := testMoveType()
  moveType.ID = unitTestTestID
  return moveType
}

func UnitTestTerrainType() *types.TerrainType {
  terrainType := testTerrainType()
  terrainType.ID = unitTestTestID
  return terrainType
}

func UnitTestItemType() *types.ItemType {
  itemType := testItemType()
  itemType.ID = unitTestTestID
  return itemType
}

func UnitTestGameType() *models.GameType {
  gameType := testGameType(unitTestTestID, unitTestTestID, unitTestTestID)
  gameType.ID = unitTestTestID
  return gameType
}

func testGameType(botTypeID, terrainTypeID, itemTypeID int64) *models.GameType {
  return &models.GameType{
		Name:         "testGameType",
		Version:      "testGameType",
		BotTypes:     []int64{botTypeID},
		TerrainTypes: []int64{terrainTypeID},
		ItemTypes:    []int64{itemTypeID},
		ApiFuncs:     []string{"me", "canMove", "move"},
		NumTeams:     1,
		CreatorID:    12345,
		Description:  "testGameType",
	}
}

func testBotType(moveTypeID, attackTypeID int64) *types.BotType {
  return &types.BotType{
		Name:              "testBot",
		AttackTypeIDs:     []int64{attackTypeID},
		MoveTypeIDs:       []int64{moveTypeID},
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
}

func testAttackType() *types.AttackType {
  return &types.AttackType{
		Name:             "testAttackType",
		Damage:           1,
		Delay:            0,
		Range:            0,
		Accuracy:         1,
		AttackDelayDealt: 0,
		MoveDelayDealt:   0,
	}
}

func testMoveType() *types.MoveType {
  return &types.MoveType{
		Name:  "testMoveType",
		Delay: 0,
		TakesDelayFromTerrain: false,
	}
}

func testTerrainType() *types.TerrainType {
  return &types.TerrainType{
		Name:            "testTerrain",
		CanBeOccupied:   true,
		MoveDelayFactor: 1,
		DamagePenalty:   1,
	}
}

func testItemType() *types.ItemType {
  return &types.ItemType{
		Name:              "testItem",
		MoveDelayFactor:   1,
		DamageFactor:      1,
		AttackDelayFactor: 1,
		RangeFactor:       1,
		AccuracyFactor:    1,
		SpawnDelayFactor:  1,
	}
}

func testMap(gameTypeID int64) *models.Map {
  return &models.Map{
		Name:       "testMap",
		GameTypeID: gameTypeID,
		RoundLimit: 10,
	}
}
