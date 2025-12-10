# Анализ системы размещения технологий на сетке

## 🔍 Текущее состояние

### Что работает правильно

✅ **TechParser уже поддерживает переменные:**
- Собирает переменные из начала файла (`@1940 = 10`)
- Резолвит переменные в координатах (`x = @RADAR` → `x = 3`)
- Метод `resolveVariable()` корректно обрабатывает ссылки

✅ **Парсинг координат:**
- Метод `parsePosition()` извлекает X и Y
- Поддерживает как числа, так и идентификаторы (переменные)
- Сохраняет в `domain.Position{X, Y}`

### ❌ Что НЕ работает

**Проблема 1: Потеря информации о переменных**
```go
// Текущий код в parsePosition():
resolvedValue := tp.resolveVariable(rawValue)  // @RADAR → "3"
val, err := strconv.Atoi(resolvedValue)        // "3" → 3
pos.X = value  // Сохраняем только число 3, теряем "@RADAR"
```

**Результат:** Мы знаем, что X=3, но не знаем, что это было `@RADAR`.

**Проблема 2: Отсутствие хранения переменных**
- `domain.Position` содержит только `X int, Y int`
- Нет полей для хранения имён переменных (`@RADAR`, `@1940`)
- Невозможно восстановить исходное представление

**Проблема 3: Отрицательные координаты**
- Radio Communications использует X от -4 до 0
- Текущий код должен работать, но нужно проверить

---

## 📋 План доработок

### Этап 1: Расширение domain.Position

**Цель:** Сохранять как числовые значения, так и имена переменных.

**Изменения в `internal/domain/position.go`:**

```go
type Position struct {
    X     int    // Числовое значение X-координаты
    Y     int    // Числовое значение Y-координаты
    XVar  string // Имя переменной для X (например, "@RADAR")
    YVar  string // Имя переменной для Y (например, "@1940")
}
```

**Преимущества:**
- Сохраняем оригинальные имена переменных
- Можем показывать в UI: "X: @RADAR (3)" вместо просто "X: 3"
- Упрощаем редактирование (snap to variables)

---

### Этап 2: Обновление TechParser

**Цель:** Сохранять имена переменных при парсинге.

**Изменения в `internal/parser/tech_parser.go`:**

```go
func (tp *TechParser) parsePosition(block *BlockStatement) domain.Position {
    pos := domain.Position{}
    
    for _, stmt := range block.Statements {
        assignStmt, ok := stmt.(*AssignmentStatement)
        if !ok {
            continue
        }
        
        var rawValue string
        var varName string  // NEW: сохраняем имя переменной
        
        switch v := assignStmt.Value.(type) {
        case *NumberLiteral:
            rawValue = v.Value
            varName = ""  // Прямое число, нет переменной
        case *Identifier:
            rawValue = v.Value
            if strings.HasPrefix(v.Value, "@") {
                varName = v.Value  // NEW: сохраняем @RADAR
            }
        }
        
        // Резолвим для получения числового значения
        resolvedValue := tp.resolveVariable(rawValue)
        val, err := strconv.Atoi(resolvedValue)
        if err != nil {
            continue
        }
        
        switch assignStmt.Name.Value {
        case "x":
            pos.X = val
            pos.XVar = varName  // NEW: сохраняем переменную
        case "y":
            pos.Y = val
            pos.YVar = varName  // NEW: сохраняем переменную
        }
    }
    
    return pos
}
```

---

### Этап 3: Экспорт словаря переменных

**Цель:** Предоставить доступ к словарю переменных для UI и других компонентов.

**Изменения в `internal/parser/tech_parser.go`:**

```go
// Добавить метод для получения переменных
func (tp *TechParser) GetVariables() map[string]string {
    return tp.variables
}

// Добавить методы для получения типизированных словарей
func (tp *TechParser) GetHorizontalVariables() map[string]int {
    result := make(map[string]int)
    for key, value := range tp.variables {
        if !strings.HasPrefix(key, "@") {
            continue
        }
        // Проверяем, что это горизонтальная переменная
        // (не содержит цифры года в имени)
        if !containsYear(key) {
            if val, err := strconv.Atoi(value); err == nil {
                result[key] = val
            }
        }
    }
    return result
}

func (tp *TechParser) GetVerticalVariables() map[string]int {
    result := make(map[string]int)
    for key, value := range tp.variables {
        if !strings.HasPrefix(key, "@") {
            continue
        }
        // Проверяем, что это вертикальная переменная (год)
        if containsYear(key) {
            if val, err := strconv.Atoi(value); err == nil {
                result[key] = val
            }
        }
    }
    return result
}

func containsYear(s string) bool {
    // Проверяем наличие 4-значного года в строке
    for i := 0; i < len(s)-3; i++ {
        if s[i] >= '0' && s[i] <= '9' &&
           s[i+1] >= '0' && s[i+1] <= '9' &&
           s[i+2] >= '0' && s[i+2] <= '9' &&
           s[i+3] >= '0' && s[i+3] <= '9' {
            return true
        }
    }
    return false
}
```

---

### Этап 4: Обновление TechnologyLoader

**Цель:** Передавать словарь переменных вместе с технологиями.

**Изменения в `internal/app/technology_loader.go`:**

```go
type TechnologyData struct {
    Technologies []*domain.Technology
    Variables    map[string]string  // NEW: словарь переменных
}

func (tl *TechnologyLoader) LoadAllTechnologiesWithVars() (*TechnologyData, error) {
    // ... существующий код загрузки ...
    
    // Получаем переменные из последнего парсера
    var variables map[string]string
    if len(allTechs) > 0 {
        // Переменные одинаковые для всех технологий в файле
        // Можем взять из любого парсера
        variables = techParser.GetVariables()
    }
    
    return &TechnologyData{
        Technologies: allTechs,
        Variables:    variables,
    }, nil
}
```

---

### Этап 5: Обновление UI для отображения

**Цель:** Показывать переменные в интерфейсе.

**Изменения в `internal/ui/scenes/tech_viewer.go`:**

```go
// При отрисовке узла технологии
func (s *TechViewerScene) drawTechNode(screen *ebiten.Image, tech *domain.Technology) {
    // ... существующий код отрисовки ...
    
    // Показываем координаты с переменными
    coordText := fmt.Sprintf("(%s, %s)", 
        formatCoordinate(tech.Position.X, tech.Position.XVar),
        formatCoordinate(tech.Position.Y, tech.Position.YVar))
    
    ebitenutil.DebugPrintAt(screen, coordText, x, y+20)
}

func formatCoordinate(value int, varName string) string {
    if varName != "" {
        return fmt.Sprintf("%s=%d", varName, value)
    }
    return fmt.Sprintf("%d", value)
}
```

---

### Этап 6: Улучшение DetectSubTrees

**Цель:** Группировать технологии по X-переменным, а не по числовым значениям.

**ВАЖНО:** Диапазоны НЕ фиксированные! Они определяются переменными в каждом файле технологий.

**Правильный алгоритм:**

1. **Группировка по XVar** (имени переменной):
   - Все технологии с `XVar = "@RADAR"` → одна группа
   - Все технологии с `XVar = "@HQ"` → другая группа
   - Технологии без переменной (прямые числа) → отдельные группы

2. **Определение границ sub-tree**:
   - Сортируем уникальные X-значения
   - Находим разрывы (gap > 5)
   - Каждый непрерывный диапазон = sub-tree

3. **Именование sub-tree**:
   - По категориям технологий в группе
   - Или по диапазону переменных

**Изменения в `internal/app/technology_loader.go`:**

```go
func (tl *TechnologyLoader) DetectSubTrees(
    folderName string, 
    technologies []*domain.Technology,
) []*domain.SubTree {
    if len(technologies) == 0 {
        return nil
    }
    
    // Шаг 1: Группируем по XVar (переменной)
    varGroups := make(map[string][]*domain.Technology)
    uniqueXValues := make(map[int]bool)
    
    for _, tech := range technologies {
        // Используем XVar как ключ группировки
        groupKey := tech.Position.XVar
        if groupKey == "" {
            // Если нет переменной, используем числовое значение
            groupKey = fmt.Sprintf("X%d", tech.Position.X)
        }
        
        varGroups[groupKey] = append(varGroups[groupKey], tech)
        uniqueXValues[tech.Position.X] = true
    }
    
    // Шаг 2: Сортируем уникальные X-значения
    xValues := make([]int, 0, len(uniqueXValues))
    for x := range uniqueXValues {
        xValues = append(xValues, x)
    }
    sort.Ints(xValues)
    
    // Шаг 3: Находим разрывы и создаём sub-trees
    subTrees := make([]*domain.SubTree, 0)
    currentRange := []int{xValues[0]}
    
    for i := 1; i < len(xValues); i++ {
        gap := xValues[i] - xValues[i-1]
        
        if gap > 5 {
            // Разрыв найден - создаём sub-tree для текущего диапазона
            subTree := createSubTreeForRange(
                currentRange[0], 
                currentRange[len(currentRange)-1],
                technologies,
                folderName,
            )
            subTrees = append(subTrees, subTree)
            
            // Начинаем новый диапазон
            currentRange = []int{xValues[i]}
        } else {
            currentRange = append(currentRange, xValues[i])
        }
    }
    
    // Добавляем последний диапазон
    if len(currentRange) > 0 {
        subTree := createSubTreeForRange(
            currentRange[0],
            currentRange[len(currentRange)-1],
            technologies,
            folderName,
        )
        subTrees = append(subTrees, subTree)
    }
    
    return subTrees
}

func createSubTreeForRange(
    xMin, xMax int,
    allTechs []*domain.Technology,
    folderName string,
) *domain.SubTree {
    // Фильтруем технологии в диапазоне
    techs := make([]*domain.Technology, 0)
    categorySet := make(map[string]bool)
    
    for _, tech := range allTechs {
        if tech.Position.X >= xMin && tech.Position.X <= xMax {
            techs = append(techs, tech)
            for _, cat := range tech.Categories {
                categorySet[cat] = true
            }
        }
    }
    
    categories := mapKeysToSlice(categorySet)
    name := identifySubTreeName(categorySet, folderName)
    
    return &domain.SubTree{
        Name:         name,
        XMin:         xMin,
        XMax:         xMax,
        Technologies: techs,
        Categories:   categories,
    }
}
```

**Преимущества этого подхода:**
- ✅ Работает с любыми переменными из файла
- ✅ Не зависит от хардкода диапазонов
- ✅ Автоматически адаптируется к разным модам
- ✅ Поддерживает отрицательные координаты

---

## 🎯 Приоритеты реализации

### HIGH Priority (критично для корректного отображения)

1. **Этап 1: Расширение domain.Position** ⭐⭐⭐
   - Добавить поля XVar, YVar
   - Обновить NewPosition()
   - **Время:** 15 минут

2. **Этап 2: Обновление TechParser** ⭐⭐⭐
   - Сохранять имена переменных в parsePosition()
   - **Время:** 30 минут

3. **Этап 6: Улучшение DetectSubTrees** ⭐⭐⭐
   - Использовать фиксированные диапазоны X для sub-trees
   - **Время:** 45 минут

### MEDIUM Priority (улучшение UX)

4. **Этап 3: Экспорт словаря переменных** ⭐⭐
   - GetVariables(), GetHorizontalVariables(), GetVerticalVariables()
   - **Время:** 30 минут

5. **Этап 5: Обновление UI** ⭐⭐
   - Показывать переменные в координатах
   - **Время:** 20 минут

### LOW Priority (для будущего редактора)

6. **Этап 4: TechnologyLoader с переменными** ⭐
   - Структура TechnologyData
   - **Время:** 20 минут

---

## 🧪 План тестирования

### Тест 1: Support Folder (простой)
- ✅ Должен работать сейчас (нет переменных или простые)
- Проверить: все технологии на своих местах

### Тест 2: Electronics Folder (сложный)
- ❌ Сейчас не работает
- После доработки проверить:
  - Radio Communications (-4..0) отдельно
  - Electronics & Radar (1..5) отдельно
  - Computing (7..9) отдельно
  - Rockets (12..16) отдельно
  - Jets (21..23) отдельно
  - Nuclear (29..33) отдельно

### Тест 3: Отрицательные координаты
- Проверить: Radio Communications корректно отображается слева

---

## 📊 Ожидаемый результат

**До:**
```
electronics_folder: все технологии в одном sub-tree
X: 0, 1, 2, 3, 4, 5, 7, 8, 9, 12, 13, 14, ...
```

**После:**
```
electronics_folder:
  ├─ Radio Communications (X: -4..-1)
  │  ├─ HQ_communications (X=-2, Y=2)
  │  ├─ radio_technology (X=-2, Y=4)
  │  └─ infantry_radio (X=-1, Y=6)
  │
  ├─ Electronics & Radar (X: 1..5)
  │  ├─ electronic_mechanical_engineering (X=1, Y=0)
  │  ├─ radio_detection (X=2, Y=4)
  │  └─ early_radar (X=3, Y=8)
  │
  ├─ Computing & Encryption (X: 7..9)
  │  ├─ mechanical_computing (X=8, Y=2)
  │  └─ electronic_computing_machine (X=7, Y=6)
  │
  ├─ Rockets & Missiles (X: 12..16)
  │  ├─ rocket_engines (X=13, Y=10)
  │  └─ ballistic_missiles (X=16, Y=20)
  │
  ├─ Jets (X: 21..23)
  │  ├─ jet_engine_theory (X=22, Y=10)
  │  └─ jet_engines (X=23, Y=16)
  │
  └─ Nuclear (X: 29..33)
     ├─ nuclear_reactor (X=31, Y=16)
     └─ nuclear_bomb (X=31, Y=20)
```

---

## 🚀 Следующие шаги

1. Начать с **Этапа 1** (domain.Position)
2. Затем **Этап 2** (TechParser)
3. Затем **Этап 6** (DetectSubTrees с фиксированными диапазонами)
4. Протестировать на electronics_folder
5. Если работает - доделать остальные этапы

**Готов начать реализацию?**
