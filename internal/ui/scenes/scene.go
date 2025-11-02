package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/shinomontaz/hoi4_visual_modder/internal/app"
)

// Scene represents a screen/state in the application
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	OnEnter()
	OnExit()
}

// SceneType identifies different scene types
type SceneType int

const (
	SceneStartup SceneType = iota
	SceneFileViewer
	SceneFocusEditor
	SceneTechEditor
)

// SceneManager manages scene transitions
type SceneManager struct {
	currentScene Scene
	scenes       map[SceneType]Scene
	state        *app.State
}

// NewSceneManager creates a new SceneManager
func NewSceneManager(state *app.State) *SceneManager {
	sm := &SceneManager{
		scenes: make(map[SceneType]Scene),
		state:  state,
	}
	
	// Register scenes
	sm.scenes[SceneStartup] = NewStartupScene(sm, state)
	sm.scenes[SceneFileViewer] = NewFileViewerScene(sm, state)
	// TODO: Add other scenes when implemented
	
	// Start with startup scene
	sm.SwitchTo(SceneStartup)
	
	return sm
}

// Update updates the current scene
func (sm *SceneManager) Update() error {
	if sm.currentScene != nil {
		return sm.currentScene.Update()
	}
	return nil
}

// Draw renders the current scene
func (sm *SceneManager) Draw(screen *ebiten.Image) {
	if sm.currentScene != nil {
		sm.currentScene.Draw(screen)
	}
}

// SwitchTo switches to a different scene
func (sm *SceneManager) SwitchTo(sceneType SceneType) {
	if sm.currentScene != nil {
		sm.currentScene.OnExit()
	}
	
	if scene, exists := sm.scenes[sceneType]; exists {
		sm.currentScene = scene
		sm.currentScene.OnEnter()
	}
}

// GetCurrentScene returns the current scene type
func (sm *SceneManager) GetCurrentScene() Scene {
	return sm.currentScene
}
