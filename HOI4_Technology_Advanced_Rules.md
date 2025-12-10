# Правила работы с технологиями HOI4 - Дополнение

## 1. Overlay Folders и логика блокировки технологий

### Что такое Overlay Folders

Overlay folders — это специальные "наложения" поверх обычных technology folders, которые блокируют доступ к технологиям до выполнения определенных условий.

### Структура в файле common/technology_tags/*.txt

```txt
technology_folders = {
    # Основная папка технологий
    infantry_folder = {
        ledger = army
    }

    # Overlay для блокировки
    infantry_overlay_folder = {
        available = {
            NOT = { has_country_flag = UNLOCK:infantry_folder }
        }
    }
}
```

### Механизм работы

1. **Основной folder** (например, `infantry_folder`) — содержит сами технологии
2. **Overlay folder** (например, `infantry_overlay_folder`) — показывается поверх основного и блокирует доступ
3. **Условие разблокировки** — через `available = { ... }` определяется, когда overlay исчезает

### Логика overlay

- Если условие в `available` **НЕ выполнено** — overlay показывается, технологии заблокированы
- Если условие **выполнено** — overlay исчезает, технологии становятся доступны
- Обычно используется отрицание `NOT = { has_country_flag = UNLOCK:folder_name }`

### Примеры из файла

```txt
# Пример 1: Простая блокировка через флаг
support_overlay_folder = { 
    available = {
        NOT = { has_country_flag = UNLOCK:support_folder }
    }
}

# Пример 2: Множественные условия
air_techs_overlay_folder = {
    available = {
        NOT = { has_country_flag = UNLOCK:air_techs_folder }
        NOT = { has_country_flag = UNLOCK:air_techs_folder_army }
        NOT = { has_country_flag = UNLOCK:air_techs_folder_navy }
    }
}

# Пример 3: Блокировка electronics и nuclear одновременно
electronics_overlay_folder = { 
    available = { 
        NOT = { has_country_flag = UNLOCK:electronics_folder }
        NOT = { has_country_flag = UNLOCK:nuclear_folder }
    }
}
```

### Country-specific технологии через available

Обычные technology folders могут иметь условия доступности для конкретных стран:

```txt
# Только для стран с флагом coastal_state
naval_folder = {
    ledger = navy
    available = {
        has_country_flag = coastal_state
        NOT = { has_dlc = "Man the Guns" }
    }
}

# Только для Германии (авиация)
luftwaffe_folder = { 
    ledger = air
    available = { 
        has_country_flag = GER_air
    }
}

# Только для стран с определённым DLC и флагом
mtgnavalfolder = {
    ledger = navy
    available = {
        has_country_flag = coastal_state_MTG
    }
}
```

### Алгоритм для инструмента

```python
def is_folder_available_for_country(folder_data, country_tag, country_flags, dlcs):
    # Проверяет доступность folder для страны
    # folder_data: словарь с данными folder из technology_tags
    # country_tag: тег страны (например, "GER")
    # country_flags: список флагов страны
    # dlcs: список установленных DLC

    if "available" not in folder_data:
        return True  # Нет условий = доступно всем

    conditions = folder_data["available"]

    # Проверяем каждое условие
    for condition_type, condition_value in conditions.items():
        if condition_type == "has_country_flag":
            if condition_value not in country_flags:
                return False

        elif condition_type == "NOT":
            # Обработка отрицания
            for neg_condition_type, neg_condition_value in condition_value.items():
                if neg_condition_type == "has_country_flag":
                    if neg_condition_value in country_flags:
                        return False
                elif neg_condition_type == "has_dlc":
                    if neg_condition_value in dlcs:
                        return False

    return True


def get_overlay_folder_for(base_folder_name):
    # Получает имя overlay folder для базового folder
    # Правило: base_folder_name + "_overlay"
    return f"{base_folder_name}_overlay_folder"


def is_folder_unlocked(base_folder_name, country_flags):
    # Проверяет, разблокирован ли folder для страны
    unlock_flag = f"UNLOCK:{base_folder_name}"
    return unlock_flag in country_flags
```

### Типичные паттерны overlay

| Паттерн | Описание | Пример |
|---------|----------|--------|
| `NOT = { has_country_flag = UNLOCK:xxx_folder }` | Блокировка до получения флага через фокус/событие | infantry_overlay_folder |
| Несколько `NOT` условий | Разблокировка при любом из флагов | air_techs_overlay_folder |
| `has_country_flag = specific_flag` | Доступно только странам с флагом | luftwaffe_folder (GER_air) |
| `has_dlc = "DLC_Name"` | Зависимость от DLC | naval_folder |
| `major_country = yes` | Только для major стран | Комбинируется с другими |

---

## 2. Где искать has_country_flag = UNLOCK:naval_doctrine_folder

### Источники флагов разблокировки

Флаги типа `UNLOCK:folder_name` устанавливаются через:

1. **National Focus** (Национальные фокусы) — файлы в `common/national_focus/*.txt`
2. **Events** (События) — файлы в `events/*.txt`
3. **Decisions** (Решения) — файлы в `common/decisions/*.txt`
4. **On actions** — файлы в `common/on_actions/*.txt`
5. **History files** (Стартовые настройки) — файлы в `history/countries/*.txt`

### Поиск в National Focus

В фокусах флаги устанавливаются через блок `completion_reward`:

```txt
focus = {
    id = unlock_naval_doctrine
    ...
    completion_reward = {
        set_country_flag = UNLOCK:naval_doctrine_folder
    }
}
```

### Поиск в Events

В событиях флаги устанавливаются через `immediate` или `option`:

```txt
country_event = {
    id = tech.1
    ...
    immediate = {
        set_country_flag = UNLOCK:naval_doctrine_folder
    }

    option = {
        name = tech.1.a
        set_country_flag = UNLOCK:naval_doctrine_folder
    }
}
```

### Поиск в History files

В истории стран флаги могут быть установлены изначально:

```txt
# history/countries/GER - Germany.txt
set_country_flag = UNLOCK:infantry_folder
set_country_flag = UNLOCK:tank_techs_folder
set_country_flag = GER_air
```

### Алгоритм поиска для инструмента

```python
def find_unlock_sources(folder_name, game_path, mod_path):
    # Найти все источники флага разблокировки
    # folder_name: имя folder (например, "naval_doctrine_folder")
    # Returns: список источников с типом и путем к файлу

    unlock_flag = f"UNLOCK:{folder_name}"
    sources = []

    # 1. Поиск в National Focus
    focus_files = glob(f"{game_path}/common/national_focus/*.txt")
    focus_files += glob(f"{mod_path}/common/national_focus/*.txt")

    for file in focus_files:
        content = read_file(file)
        if unlock_flag in content:
            sources.append({
                "type": "national_focus",
                "file": file,
                "flag": unlock_flag
            })

    # 2. Поиск в Events
    event_files = glob(f"{game_path}/events/*.txt")
    event_files += glob(f"{mod_path}/events/*.txt")

    for file in event_files:
        content = read_file(file)
        if unlock_flag in content:
            sources.append({
                "type": "event",
                "file": file,
                "flag": unlock_flag
            })

    # 3. Поиск в Decisions
    decision_files = glob(f"{game_path}/common/decisions/*.txt")
    decision_files += glob(f"{mod_path}/common/decisions/*.txt")

    for file in decision_files:
        content = read_file(file)
        if unlock_flag in content:
            sources.append({
                "type": "decision",
                "file": file,
                "flag": unlock_flag
            })

    # 4. Поиск в History
    history_files = glob(f"{game_path}/history/countries/*.txt")
    history_files += glob(f"{mod_path}/history/countries/*.txt")

    for file in history_files:
        content = read_file(file)
        if unlock_flag in content:
            sources.append({
                "type": "history",
                "file": file,
                "flag": unlock_flag
            })

    return sources
```

### Практический пример

Для `UNLOCK:naval_doctrine_folder` нужно искать:
- В фокусах морских держав (UK, USA, JAP и т.д.)
- В событиях, связанных с военной реформой
- В стартовых файлах истории major стран

---

## 3. Несколько деревьев в одном folder (Sub-technologies)

### Проблема: Engineering folder

В игре вкладка "Engineering and advanced tech" (`electronics_folder`) содержит **4 независимых поддерева**:

1. **Electronic engineering** (электроника, компьютеры, радар)
2. **Experimental rockets** (ракетные технологии)
3. **Jets & Aircraft engines** (реактивные двигатели)
4. **Atomic research** (атомные исследования)

### Как это работает

Каждое поддерево начинается в **своих координатах** внутри общего folder. Это достигается через:

1. **Разные стартовые позиции** — каждое дерево имеет свою точку (0,0)
2. **Независимые ветки** — технологии не связаны dependencies между деревьями
3. **Общий folder** — все технологии принадлежат одному `electronics_folder`

### Структура в technologies/*.txt

```txt
# Дерево 1: Electronic engineering (стартовая позиция x=0)
electronic_mechanical_engineering = {
    folder = {
        name = electronics_folder
        position = { x = 0 y = 0 }  # Начало первого дерева
    }
    categories = { electronics computing_tech }
}

# Следующая технология в том же дереве
radio = {
    folder = {
        name = electronics_folder
        position = { x = 0 y = 2 }  # x=0 = первое дерево
    }
    dependencies = {
        electronic_mechanical_engineering = 1
    }
}

# Дерево 2: Experimental Rockets (стартовая позиция x=6)
rocket_engines = {
    folder = {
        name = electronics_folder
        position = { x = 6 y = 0 }  # Начало второго дерева
    }
    categories = { rocketry }
}

# Дерево 3: Jets (стартовая позиция x=12)
jet_engines = {
    folder = {
        name = electronics_folder
        position = { x = 12 y = 0 }  # Начало третьего дерева
    }
    categories = { jet_engine }
}

# Дерево 4: Atomic (стартовая позиция x=18)
nuclear_reactor = {
    folder = {
        name = electronics_folder
        position = { x = 18 y = 0 }  # Начало четвертого дерева
    }
    categories = { nuclear }
}
```

### Координаты поддеревьев

Поддеревья **группируются по X-координате**:

| Поддерево | Стартовая X | Категории | Описание |
|-----------|-------------|-----------|----------|
| Electronic engineering | 0-5 | electronics, computing_tech, radar_tech, radio_tech | Электроника, радар, компьютеры |
| Experimental Rockets | 6-11 | rocketry, mot_rockets | Ракеты |
| Jets & Aircraft engines | 12-17 | jet_technology, jet_engine | Реактивные двигатели |
| Atomic research | 18-23 | nuclear | Атомные исследования |

### Правила координат

1. **X-координата** определяет, к какому поддереву относится технология
2. **Y-координата** — вертикальная позиция в дереве (чем больше Y, тем ниже)
3. **Отступ между деревьями** — обычно 6 единиц по X (0, 6, 12, 18...)
4. **Независимость** — деревья не пересекаются по X

### Алгоритм определения поддеревьев

```python
def detect_sub_trees(folder_name, technologies):
    # Определяет поддеревья внутри folder
    # folder_name: имя folder (например, "electronics_folder")
    # technologies: список всех технологий
    # Returns: список поддеревьев с диапазонами координат

    # Отфильтровать технологии для данного folder
    folder_techs = []
    for tech_id, tech_data in technologies.items():
        for folder_block in tech_data.get("folder", []):
            if folder_block.get("name") == folder_name:
                position = folder_block.get("position", {})
                x = position.get("x", 0)
                y = position.get("y", 0)
                folder_techs.append({
                    "id": tech_id,
                    "x": x,
                    "y": y,
                    "categories": tech_data.get("categories", [])
                })

    # Группировка по X-координате (кластеризация)
    # Предполагаем отступ в 6 единиц
    sub_trees = []
    x_values = sorted(set([tech["x"] for tech in folder_techs]))

    if not x_values:
        return []

    # Группируем по диапазонам X
    current_tree = {
        "x_min": x_values[0],
        "x_max": x_values[0],
        "technologies": []
    }

    for tech in folder_techs:
        # Если X далеко от текущего дерева - начать новое
        if tech["x"] - current_tree["x_max"] > 5:
            sub_trees.append(current_tree)
            current_tree = {
                "x_min": tech["x"],
                "x_max": tech["x"],
                "technologies": []
            }
        else:
            current_tree["x_max"] = max(current_tree["x_max"], tech["x"])

        current_tree["technologies"].append(tech)

    # Добавить последнее дерево
    if current_tree["technologies"]:
        sub_trees.append(current_tree)

    # Определить категории для каждого поддерева
    for tree in sub_trees:
        categories = set()
        for tech in tree["technologies"]:
            categories.update(tech["categories"])
        tree["categories"] = list(categories)

    return sub_trees


def identify_sub_tree_by_categories(sub_trees):
    # Идентифицировать поддеревья по категориям
    # sub_trees: результат detect_sub_trees
    # Returns: sub_trees с добавленными именами

    category_to_name = {
        "electronics": "Electronic Engineering",
        "computing_tech": "Electronic Engineering",
        "radar_tech": "Electronic Engineering",
        "rocketry": "Experimental Rockets",
        "mot_rockets": "Experimental Rockets",
        "jet_technology": "Jets & Aircraft Engines",
        "jet_engine": "Jets & Aircraft Engines",
        "nuclear": "Atomic Research"
    }

    for tree in sub_trees:
        # Найти имя по категориям
        for category in tree["categories"]:
            if category in category_to_name:
                tree["name"] = category_to_name[category]
                break

        if "name" not in tree:
            tree["name"] = f"Sub-tree {tree['x_min']}-{tree['x_max']}"

    return sub_trees
```

### Визуализация в инструменте

```python
# Пример использования
technologies = load_all_technologies(game_path, mod_path)
sub_trees = detect_sub_trees("electronics_folder", technologies)
sub_trees = identify_sub_tree_by_categories(sub_trees)

# Результат:
# [
#   {
#     "name": "Electronic Engineering",
#     "x_min": 0,
#     "x_max": 5,
#     "categories": ["electronics", "computing_tech", "radar_tech"],
#     "technologies": [...]
#   },
#   {
#     "name": "Experimental Rockets",
#     "x_min": 6,
#     "x_max": 11,
#     "categories": ["rocketry"],
#     "technologies": [...]
#   },
#   ...
# ]

# В UI показывать как отдельные вкладки или колонки
for tree in sub_trees:
    render_sub_tree(tree["name"], tree["technologies"])
```

---

## 4. Координаты и размещение поддеревьев

### Система координат технологий

В HOI4 для технологий используется **сеточная система координат**:

- **X** — горизонтальная позиция (слева направо)
- **Y** — вертикальная позиция (сверху вниз)
- **Единица измерения** — не пиксели, а слоты сетки

### Отличие от национальных фокусов

| Аспект | National Focus | Technologies |
|--------|----------------|--------------|
| Единица измерения | Фиксированная сетка | Настраиваемый slotsize в GUI |
| Координата (0,0) | Центр дерева | Верхний левый угол gridbox |
| Отрицательные координаты | Возможны | Обычно не используются |

### Slotsize и gridbox

Размер слота определяется в GUI файлах (`interface/countrytechtreeview.gui`):

```
gridboxType = {
    name = "technology_folder_electronics"
    position = { x = 100 y = 100 }  # Позиция gridbox в окне (в пикселях)
    slotsize = { width = 72 height = 72 }  # Размер одного слота (в пикселях)
    max_slots_horizontal = 24  # Максимум слотов по горизонтали
}
```

**Важно:**
- `position` в gridbox — позиция в **пикселях** относительно окна
- `position` в технологии — позиция в **слотах** относительно gridbox
- Реальная позиция технологии = gridbox.position + (tech.position × slotsize)

### Правила размещения поддеревьев

1. **Начало каждого дерева** — Y=0 (верхняя граница)
2. **Отступ между деревьями** — обычно 6 слотов по X
3. **Вертикальное развитие** — увеличение Y для зависимостей
4. **Связи (dependencies)** — только внутри одного поддерева (одинаковый диапазон X)

### Примеры координат

```txt
# Electronic Engineering (X: 0-5)
electronic_mechanical_engineering:  x=0, y=0  # Стартовая технология
radio:                              x=0, y=2
radio_detection:                    x=2, y=4
improved_radar:                     x=2, y=6

# Experimental Rockets (X: 6-11)
rocket_engines:                     x=6, y=0  # Стартовая технология
improved_rocket_engines:            x=6, y=2
advanced_rocket_engines:            x=6, y=4

# Jets (X: 12-17)
jet_engines:                        x=12, y=0  # Стартовая технология
improved_jet_engines:               x=12, y=2

# Atomic (X: 18-23)
nuclear_reactor:                    x=18, y=0  # Стартовая технология
nuclear_bomb:                       x=18, y=4
```

### Алгоритм расчета координат для UI

```python
def calculate_pixel_position(tech_position, gridbox_config):
    # Конвертирует координаты технологии в пиксели
    # tech_position: {"x": int, "y": int} из определения технологии
    # gridbox_config: конфигурация gridbox из GUI
    # Returns: {"x": int, "y": int} в пикселях

    slot_width = gridbox_config["slotsize"]["width"]
    slot_height = gridbox_config["slotsize"]["height"]

    gridbox_x = gridbox_config["position"]["x"]
    gridbox_y = gridbox_config["position"]["y"]

    pixel_x = gridbox_x + (tech_position["x"] * slot_width)
    pixel_y = gridbox_y + (tech_position["y"] * slot_height)

    return {"x": pixel_x, "y": pixel_y}


def draw_technology_tree(folder_name, sub_trees, gridbox_config):
    # Отрисовка дерева технологий с поддеревьями

    for sub_tree in sub_trees:
        # Заголовок поддерева
        header_x = (sub_tree["x_min"] + sub_tree["x_max"]) / 2
        draw_text(sub_tree["name"], calculate_pixel_position({"x": header_x, "y": -1}, gridbox_config))

        # Технологии в поддереве
        for tech in sub_tree["technologies"]:
            pos = calculate_pixel_position(tech, gridbox_config)
            draw_technology_icon(tech["id"], pos)

            # Связи (dependencies)
            if "dependencies" in tech:
                for dep_tech_id in tech["dependencies"]:
                    dep_tech = find_technology(dep_tech_id, sub_tree["technologies"])
                    if dep_tech:
                        dep_pos = calculate_pixel_position(dep_tech, gridbox_config)
                        draw_connection_line(dep_pos, pos)
```

### Рекомендации для редактора

1. **Автоматическое определение поддеревьев** при загрузке folder
2. **Визуальное разделение** поддеревьев вертикальными линиями или цветом
3. **Независимое редактирование** каждого поддерева
4. **Валидация координат** — предупреждать если технология выходит за границы своего поддерева
5. **Snap to grid** — автоматическое выравнивание по сетке при перетаскивании

---

## Итоговая структура данных для инструмента

```json
{
  "folder_id": "electronics_folder",
  "folder_name": "Engineering and advanced tech",
  "ledger": "civilian",
  "has_overlay": true,
  "overlay_conditions": {
    "NOT": [
      {"has_country_flag": "UNLOCK:electronics_folder"},
      {"has_country_flag": "UNLOCK:nuclear_folder"}
    ]
  },
  "sub_trees": [
    {
      "name": "Electronic Engineering",
      "x_range": [0, 5],
      "categories": ["electronics", "computing_tech", "radar_tech"],
      "technologies": [
        {
          "id": "electronic_mechanical_engineering",
          "position": {"x": 0, "y": 0},
          "dependencies": [],
          "localized_name": "Electronic Mechanical Engineering"
        },
        {
          "id": "radio",
          "position": {"x": 0, "y": 2},
          "dependencies": ["electronic_mechanical_engineering"],
          "localized_name": "Radio"
        }
      ]
    },
    {
      "name": "Experimental Rockets",
      "x_range": [6, 11],
      "categories": ["rocketry"],
      "technologies": [...]
    },
    {
      "name": "Jets & Aircraft Engines",
      "x_range": [12, 17],
      "categories": ["jet_technology", "jet_engine"],
      "technologies": [...]
    },
    {
      "name": "Atomic Research",
      "x_range": [18, 23],
      "categories": ["nuclear"],
      "technologies": [...]
    }
  ]
}
```

---

## Справочные материалы

- [HOI4 Technology Modding Wiki](https://hoi4.paradoxwikis.com/Technology_modding)
- [HOI4 Engineering Technology Wiki](https://hoi4.paradoxwikis.com/Engineering_technology)
- Примеры из `common/technologies/` и `common/technology_tags/`
- GUI определения в `interface/countrytechtreeview.gui`
