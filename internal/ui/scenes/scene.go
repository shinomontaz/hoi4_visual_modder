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
	currentScene  Scene
	scenes        map[SceneType]Scene
	dynamicScenes map[string]Scene // For dynamically created scenes
	state         *app.State
}

// NewSceneManager creates a new SceneManager
func NewSceneManager(state *app.State) *SceneManager {
	sm := &SceneManager{
		scenes:        make(map[SceneType]Scene),
		dynamicScenes: make(map[string]Scene),
		state:         state,
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

// SwitchTo switches to a different scene by SceneType
func (sm *SceneManager) SwitchTo(sceneType SceneType) {
	if sm.currentScene != nil {
		sm.currentScene.OnExit()
	}
	
	if scene, exists := sm.scenes[sceneType]; exists {
		sm.currentScene = scene
		sm.currentScene.OnEnter()
	}
}

// SwitchToNamed switches to a dynamically created scene by name
func (sm *SceneManager) SwitchToNamed(name string) {
	if sm.currentScene != nil {
		sm.currentScene.OnExit()
	}
	
	if scene, exists := sm.dynamicScenes[name]; exists {
		sm.currentScene = scene
		sm.currentScene.OnEnter()
	}
}

// AddScene adds a dynamic scene
func (sm *SceneManager) AddScene(name string, scene Scene) {
	sm.dynamicScenes[name] = scene
}

// GetCurrentScene returns the current scene type
func (sm *SceneManager) GetCurrentScene() Scene {
	return sm.currentScene
}
