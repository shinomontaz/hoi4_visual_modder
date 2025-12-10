# Правила работы с технологиями в Hearts of Iron 4

## 1. Определение списка деревьев технологий (Technology Folders)

### Где искать информацию
Деревья технологий определяются через **folders** (папки/вкладки) в файлах:
- `common/technology_tags/*.txt` - определение списка всех folders
- `common/technologies/*.txt` - сами технологии с указанием принадлежности к folders

### Структура technology_tags

В файлах `common/technology_tags/*.txt` находится блок `technology_folders`:

```txt
technology_folders = {
    infantry_folder
    support_folder
    armor_folder
    artillery_folder
    naval_folder
    air_folder
    engineering_folder
    industry_folder
    ...
}
```

## 3. Локализация названий веток технологий (technology folders)

Названия вкладок (веток) технологий не задаются в `common/technology_tags/`, а берутся из локализации по ключу `<folder_id>_name`. Например, для `infantry_folder` используется ключ `infantry_folder_name`. Локализация лежит в `localisation/<language>/*_l_<language>.yml` как в базе игры, так и в моде. При конфликте значение из мода должно перекрывать базовое.

Пример (english):

```
l_english:
infantry_folder_name:0 "Infantry"
support_folder_name:0 "Support"
armor_folder_name:0 "Armor"
artillery_folder_name:0 "Artillery"
naval_folder_name:0 "Naval"
air_folder_name:0 "Air"
engineering_folder_name:0 "Engineering"
industry_folder_name:0 "Industry"
```

#### Алгоритм для инструмента:

1. Есть `folder_id` (например, `infantry_folder`) из `technology_tags`.
2. Строим ключ локализации: `<folder_id> + "_name"` → `infantry_folder_name`.
3. Читаем все `*_l_<language>.yml` из `localisation/<language>/` игры, потом мода.
4. Ищем этот ключ; если найден — используем значение, если нет — показываем `folder_id` (можно обрезать `_folder` и капитализировать).


### Структура технологии

Каждая технология в `common/technologies/*.txt` содержит параметр `folder`, который указывает в каких вкладках она отображается:

```txt
my_technology = {
    folder = {
        name = infantry_folder
        position = { x = 2 y = 5 }
    }
    # Технология может быть в нескольких folders
    folder = {
        name = support_folder
        position = { x = 4 y = 3 }
    }

    categories = { infantry_weapons }

    dependencies = {
        previous_tech = 1
    }

    # Другие параметры...
}
```

### Алгоритм получения списка деревьев технологий

1. Прочитать все файлы из `common/technology_tags/` (сначала из игры, затем из мода)
2. Распарсить блок `technology_folders = { ... }`
3. Получить список всех folders (например: infantry_folder, armor_folder, naval_folder и т.д.)
4. Это и есть список всех вкладок/деревьев технологий

### Получение технологий для конкретного дерева

1. Прочитать все файлы из `common/technologies/` (сначала из игры, затем из мода с перезаписью)
2. Для каждой технологии проверить блоки `folder = { name = ... }`
3. Если `name` совпадает с нужным folder - добавить технологию в это дерево
4. Использовать `position = { x = ... y = ... }` для размещения технологии на визуальном дереве

---

## 2. Переопределение технологий для конкретной страны

### Важное правило
**Технологии НЕ переопределяются через отдельные файлы для каждой страны.** Все технологии общие для всех стран и находятся в `common/technologies/`.

### Механизмы кастомизации технологий

#### 2.1 Country-specific названия через локализацию

Технологии могут иметь уникальные названия для разных стран через файлы локализации:

**Стандартная локализация:**
```yml
my_equipment_1:0 "My Equipment"
```

**Country-specific локализация (для страны с тегом HON - Гондурас):**
```yml
HON_my_equipment_1:0 "My Honduran Equipment"
```

Игра автоматически использует локализацию с префиксом тега страны, если она существует.

#### 2.2 Условия доступности (allow и available)

Технологии могут быть ограничены для определённых стран через условия:

```txt
my_special_tech = {
    # Технология доступна только для Германии
    allow = {
        original_tag = GER
    }

    # Или для определённых идеологий
    available = {
        has_government = fascism
    }

    folder = {
        name = infantry_folder
        position = { x = 5 y = 10 }
    }

    # Остальные параметры...
}
```

#### 2.3 Ветки технологий (allow_branch)

Целые ветки технологий могут быть ограничены:

```txt
my_tech = {
    allow_branch = {
        # Эта технология и все зависимые от неё доступны только для СССР
        tag = SOV
    }

    folder = {
        name = armor_folder
        position = { x = 3 y = 7 }
    }
}
```

### Алгоритм работы с технологиями для страны

1. **Загрузить все технологии** из `common/technologies/` (базовая игра + мод)
2. **Для каждой технологии проверить условия:**
   - Есть ли блок `allow = { ... }` или `allow_branch = { ... }`
   - Если условия есть - проверить применимость к выбранной стране
   - Если условий нет - технология доступна всем странам
3. **Применить локализацию:**
   - Сначала искать локализацию с префиксом `TAG_tech_name`
   - Если нет - использовать стандартную `tech_name`
4. **Отобразить дерево технологий** с учётом фильтрации

### Переопределение в моде

Если мод хочет изменить существующую технологию:

1. Создать файл с технологией в `common/technologies/` мода
2. Использовать **то же самое имя технологии**
3. Мод **полностью перезапишет** определение технологии из базовой игры

Пример структуры:
```
mod_folder/
├── common/
│   └── technologies/
│       └── infantry.txt  # Перезаписывает infantry.txt из игры
```

### Важные замечания

- **Нельзя** создать отдельные файлы типа `technologies_GER.txt` или `technologies_SOV.txt` - это не работает
- **Правильный подход:** одно общее дерево технологий + условия доступности + country-specific локализация
- Названия файлов в `common/technologies/` могут быть любыми (infantry.txt, armor.txt, naval.txt и т.д.) - игра читает все `.txt` файлы
- Мод может добавлять **новые** технологии или **перезаписывать** существующие

---

## Практический пример для инструмента

### Чтение технологий

```python
# Псевдокод для чтения технологий

def load_technologies(game_path, mod_path, country_tag):
    technologies = {}

    # 1. Загрузить технологии из игры
    game_tech_files = glob(game_path + "/common/technologies/*.txt")
    for file in game_tech_files:
        techs = parse_paradox_file(file)
        technologies.update(techs)

    # 2. Загрузить технологии из мода (перезапись)
    mod_tech_files = glob(mod_path + "/common/technologies/*.txt")
    for file in mod_tech_files:
        techs = parse_paradox_file(file)
        technologies.update(techs)  # Перезаписывает дубликаты

    # 3. Фильтрация по стране
    filtered_techs = {}
    for tech_name, tech_data in technologies.items():
        if is_tech_available_for_country(tech_data, country_tag):
            filtered_techs[tech_name] = tech_data

    return filtered_techs

def is_tech_available_for_country(tech_data, country_tag):
    # Проверить allow и allow_branch
    if "allow" in tech_data:
        return evaluate_conditions(tech_data["allow"], country_tag)
    if "allow_branch" in tech_data:
        return evaluate_conditions(tech_data["allow_branch"], country_tag)
    # Если условий нет - доступна всем
    return True
```

### Получение folders

```python
def load_technology_folders(game_path, mod_path):
    folders = []

    # 1. Из игры
    tag_files = glob(game_path + "/common/technology_tags/*.txt")
    for file in tag_files:
        data = parse_paradox_file(file)
        if "technology_folders" in data:
            folders.extend(data["technology_folders"])

    # 2. Из мода (дополнение)
    mod_tag_files = glob(mod_path + "/common/technology_tags/*.txt")
    for file in mod_tag_files:
        data = parse_paradox_file(file)
        if "technology_folders" in data:
            folders.extend(data["technology_folders"])

    # Убрать дубликаты
    return list(set(folders))
```

---

## Итоговая структура для инструмента

### Данные для отображения

```json
{
  "folders": [
    "infantry_folder",
    "armor_folder",
    "naval_folder",
    "air_folder",
    "industry_folder"
  ],
  "technologies": {
    "infantry_weapons1": {
      "folders": [
        {"name": "infantry_folder", "position": {"x": 2, "y": 0}}
      ],
      "available_for_country": true,
      "localized_name": "GER_infantry_weapons1",
      "dependencies": ["basic_weapons"],
      "categories": ["infantry_weapons", "land"]
    }
  }
}
```

### UI инструмента

1. **Выбор folder** (вкладка дерева технологий)
2. **Отображение технологий** из этого folder с учётом:
   - Позиции (x, y)
   - Зависимостей (стрелки между технологиями)
   - Доступности для выбранной страны
3. **Редактирование технологии:**
   - Изменение позиции
   - Изменение зависимостей
   - Добавление/удаление условий allow
   - Редактирование локализации

---

## Справочные материалы

- [HOI4 Technology Modding Wiki](https://hoi4.paradoxwikis.com/Technology_modding)
- Форумы Paradox Plaza
- Примеры из базовой игры в `common/technologies/` и `common/technology_tags/`
