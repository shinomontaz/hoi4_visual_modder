# Правила размещения технологий на сетке в HOI4

## Система координат технологий

В Hearts of Iron 4 технологии размещаются на двумерной сетке с координатами **X (горизонталь)** и **Y (вертикаль)**.

### Основные принципы

1. **X-координата** определяет горизонтальное положение технологии (колонку в дереве)
2. **Y-координата** определяет вертикальное положение (строку в дереве)
3. Координаты задаются через **переменные** для удобства и единообразия
4. **Отрицательные координаты X** допустимы для левых колонок

---

## Вертикальные координаты (Y) - Временная шкала

Вертикальные позиции привязаны к **годам** и представляют хронологию исследований:

```txt
###Vertical positions
@1930 = 0
@1930_1 = 1
@1936 = 2
@1936_1 = 3
@1937 = 4
@1937_1 = 5
@1938 = 6
@1938_1 = 7
@1939 = 8
@1939_1 = 9
@1940 = 10
@1940_1 = 11
@1941 = 12
@1941_1 = 13
@1942 = 14
@1942_1 = 15
@1943 = 16
@1944 = 18
@1945 = 20
@1946 = 22
@1947 = 24
@1948 = 26
@1949 = 28
@1950 = 30
@1951 = 32
```

### Правила вертикального размещения

| Год | Y-координата | Описание |
|-----|--------------|----------|
| 1930 | 0 | Самая верхняя позиция, ранние технологии |
| 1936 | 2 | Стартовые технологии на начало игры |
| 1937-1939 | 4-8 | Ранняя война |
| 1940-1942 | 10-14 | Середина войны |
| 1943-1945 | 16-20 | Поздняя война |
| 1946-1951 | 22-32 | Послевоенные технологии |

**Важно:**
- Суффикс `_1` (например, `@1936_1 = 3`) означает промежуточную позицию между годами
- Интервал между основными годами неравномерный: 1-2 единицы в начале, затем увеличивается
- С 1944 года интервал увеличивается до 2 единиц между годами

---

## Горизонтальные координаты (X) - Функциональные группы

Горизонтальные позиции определяют **тип технологии** и её принадлежность к определённому поддереву.

### Полный список горизонтальных переменных

```txt
###Horizontal positions

## Радио и коммуникации (левая часть, отрицательные координаты)
@CORP_RADIO = -4      # Связь на уровне корпуса
@TNK_RADIO = -3       # Танковое радио
@RECON_RADIO = -2     # Радио для разведки
@HQ = -2              # Штабы и командование
@INF_RADIO = -1       # Пехотное радио
@TRCK_RADIO = 0       # Радио для транспорта

## Электроника и радар (центральная часть)
@EL_ME_ENG = 1        # Electronic Mechanical Engineering (основа)
@LISTEN = 1           # Прослушивание (listening stations)
@CONSUMER = 2         # Потребительская электроника
@RADAR = 3            # Радар
@TV = 4               # Телевидение
@AVIONICS = 5         # Авионика

## Компьютеры и шифрование (правее)
@ELEC_COMP = 7        # Электронные компьютеры
@MECH_COMP = 8        # Механические компьютеры
@BASIC_ENCRP = 9      # Базовое шифрование

## Ракеты и космос (дальше вправо)
@RCKS = 12            # Rockets (ракеты)
@RCKT_ENG = 13        # Rocket Engines (ракетные двигатели)
@MSL = 14             # Missiles (управляемые ракеты)
@AEROSPACE = 15       # Аэрокосмические технологии
@BAL_MISSILE = 16     # Ballistic Missiles (баллистические ракеты)

## Реактивные технологии (Jets)
@JET_PROTOTYPE = 21   # Прототипы реактивных двигателей
@JET_THEORY = 22      # Теория реактивных двигателей
@JET_ENG = 23         # Jet Engines (реактивные двигатели)

## Авиационные двигатели
@AIR_ENG = 25         # Air Engines (авиадвигатели)

## Атомные технологии (крайняя правая часть)
@FISS_EXP = 29        # Fission Experiments (эксперименты с делением)
@ATOM_RES = 31        # Atomic Research (атомные исследования)
@IMP_PART = 33        # Improved Particles (улучшенные частицы)
```

---

## Карта поддеревьев по X-координатам

### Визуальная схема electronics_folder

```
X:  -4  -3  -2  -1   0   1   2   3   4   5 | 7   8   9 | 12  13  14  15  16 | 21  22  23 | 25 | 29  31  33
    └──────────────────────┘ └────────────┘   └────────┘   └──────────────────┘   └────────┘   │   └─────────┘
    Radio Communications      Electronics      Computing    Rockets & Missiles       Jets       │    Nuclear
                              & Radar                                                            │
                                                                                            Air Engines
```

### Группировка поддеревьев

| X-диапазон | Название поддерева | Описание |
|------------|-------------------|----------|
| **-4 .. 0** | Radio Communications | Радиосвязь для различных родов войск |
| **1 .. 5** | Electronics & Radar | Электроника, радар, авионика, потребительская электроника |
| **7 .. 9** | Computing & Encryption | Компьютеры и криптография |
| **12 .. 16** | Rockets & Missiles | Ракетные технологии |
| **21 .. 23** | Jets | Реактивные двигатели |
| **25** | Air Engines | Авиационные двигатели |
| **29 .. 33** | Nuclear | Атомные исследования |

---

## Примеры размещения конкретных технологий

### Radio Communications (X: -4 .. 0)

```txt
electronic_mechanical_engineering = {
    folder = {
        name = electronics_folder
        position = { x = @EL_ME_ENG y = @1930 }  # x=1, y=0
    }
}

HQ_communications = {
    folder = {
        name = electronics_folder
        position = { x = @HQ y = @1936 }  # x=-2, y=2
    }
}

radio_technology = {
    folder = {
        name = electronics_folder
        position = { x = @HQ y = @1937 }  # x=-2, y=4
    }
}

infantry_radio = {
    folder = {
        name = electronics_folder
        position = { x = @INF_RADIO y = @1938 }  # x=-1, y=6
    }
}

corps_communications = {
    folder = {
        name = electronics_folder
        position = { x = @CORP_RADIO y = @1938 }  # x=-4, y=6
    }
}

tank_radio = {
    folder = {
        name = electronics_folder
        position = { x = @TNK_RADIO y = @1940 }  # x=-3, y=10
    }
}

vehicle_radio = {
    folder = {
        name = electronics_folder
        position = { x = @TRCK_RADIO y = @1939 }  # x=0, y=8
    }
}

recon_radio = {
    folder = {
        name = electronics_folder
        position = { x = @RECON_RADIO y = @1939 }  # x=-2, y=8
    }
}
```

### Electronics & Radar (X: 1 .. 5)

```txt
consumer_electronics = {
    folder = {
        name = electronics_folder
        position = { x = @TV y = @1936 }  # x=4, y=2
    }
}

radio_detection = {
    folder = {
        name = electronics_folder
        position = { x = @CONSUMER y = @1937 }  # x=2, y=4
    }
}

early_radar = {
    folder = {
        name = electronics_folder
        position = { x = @RADAR y = @1939 }  # x=3, y=8
    }
}

decimetric_radar = {
    folder = {
        name = electronics_folder
        position = { x = @RADAR y = @1940 }  # x=3, y=10
    }
}

avionics_improvements = {
    folder = {
        name = electronics_folder
        position = { x = @AVIONICS y = @1938 }  # x=5, y=6
    }
}
```

### Computing & Encryption (X: 7 .. 9)

```txt
mechanical_computing = {
    folder = {
        name = electronics_folder
        position = { x = @MECH_COMP y = @1936 }  # x=8, y=2
    }
}

electronic_computing_machine = {
    folder = {
        name = electronics_folder
        position = { x = @ELEC_COMP y = @1938 }  # x=7, y=6
    }
}

basic_encryption = {
    folder = {
        name = electronics_folder
        position = { x = @BASIC_ENCRP y = @1936 }  # x=9, y=2
    }
}
```

### Rockets & Missiles (X: 12 .. 16)

```txt
rocket_engines = {
    folder = {
        name = electronics_folder
        position = { x = @RCKT_ENG y = @1940 }  # x=13, y=10
    }
}

improved_rocket_engines = {
    folder = {
        name = electronics_folder
        position = { x = @RCKT_ENG y = @1942 }  # x=13, y=14
    }
}

ballistic_missiles = {
    folder = {
        name = electronics_folder
        position = { x = @BAL_MISSILE y = @1945 }  # x=16, y=20
    }
}
```

### Jets (X: 21 .. 23)

```txt
jet_engine_theory = {
    folder = {
        name = electronics_folder
        position = { x = @JET_THEORY y = @1940 }  # x=22, y=10
    }
}

jet_engines = {
    folder = {
        name = electronics_folder
        position = { x = @JET_ENG y = @1943 }  # x=23, y=16
    }
}
```

### Nuclear (X: 29 .. 33)

```txt
nuclear_reactor = {
    folder = {
        name = electronics_folder
        position = { x = @ATOM_RES y = @1943 }  # x=31, y=16
    }
}

nuclear_bomb = {
    folder = {
        name = electronics_folder
        position = { x = @ATOM_RES y = @1945 }  # x=31, y=20
    }
}

improved_particle_physics = {
    folder = {
        name = electronics_folder
        position = { x = @IMP_PART y = @1947 }  # x=33, y=24
    }
}
```

---

## Правила использования координат в моде

### 1. Объявление переменных

В начале файла технологий (`common/technologies/*.txt`) объявляются все переменные:

```txt
technologies = {
    ###Vertical positions
    @1930 = 0
    @1936 = 2
    @1937 = 4
    ...

    ###Horizontal positions
    @CORP_RADIO = -4
    @TNK_RADIO = -3
    ...

    # Далее идут определения технологий
}
```

### 2. Использование в технологиях

```txt
my_technology = {
    folder = {
        name = electronics_folder
        position = { x = @RADAR y = @1940 }  # Используем переменные
    }

    start_year = 1940  # Год должен соответствовать Y-координате
    ...
}
```

### 3. Создание новых технологий

При добавлении новой технологии:

1. **Определи функциональную группу** — к какому поддереву относится (радио, радар, ракеты и т.д.)
2. **Выбери X-координату** — используй существующую переменную или создай новую
3. **Выбери Y-координату** — по году доступности технологии
4. **Проверь зависимости** — технологии в одной вертикальной линии (одинаковый X) должны образовывать логическую цепочку

### 4. Добавление новых переменных

Если нужна новая горизонтальная позиция:

```txt
###Horizontal positions
...
@MY_NEW_TECH = 10  # Новая координата между BASIC_ENCRP (9) и RCKS (12)
```

Если нужна новая вертикальная позиция:

```txt
###Vertical positions
...
@1941_2 = 12.5  # Промежуточная позиция (но лучше использовать целые числа)
```

---

## Правила для инструмента редактирования

### Словари переменных для автокомплита

```python
HORIZONTAL_VARS = {
    "CORP_RADIO": -4,
    "TNK_RADIO": -3,
    "RECON_RADIO": -2,
    "HQ": -2,
    "INF_RADIO": -1,
    "TRCK_RADIO": 0,
    "EL_ME_ENG": 1,
    "CONSUMER": 2,
    "RADAR": 3,
    "TV": 4,
    "AVIONICS": 5,
    "ELEC_COMP": 7,
    "MECH_COMP": 8,
    "BASIC_ENCRP": 9,
    "RCKS": 12,
    "RCKT_ENG": 13,
    "MSL": 14,
    "AEROSPACE": 15,
    "BAL_MISSILE": 16,
    "JET_PROTOTYPE": 21,
    "JET_THEORY": 22,
    "JET_ENG": 23,
    "AIR_ENG": 25,
    "FISS_EXP": 29,
    "ATOM_RES": 31,
    "IMP_PART": 33
}

VERTICAL_VARS = {
    "1930": 0, "1930_1": 1,
    "1936": 2, "1936_1": 3,
    "1937": 4, "1937_1": 5,
    "1938": 6, "1938_1": 7,
    "1939": 8, "1939_1": 9,
    "1940": 10, "1940_1": 11,
    "1941": 12, "1941_1": 13,
    "1942": 14, "1942_1": 15,
    "1943": 16,
    "1944": 18,
    "1945": 20,
    "1946": 22,
    "1947": 24,
    "1948": 26,
    "1949": 28,
    "1950": 30,
    "1951": 32
}
```

### UI для редактора

В интерфейсе редактора показывать:
1. **Сетку с переменными** — не числа, а имена переменных (@RADAR, @1940)
2. **Подсветку поддеревьев** — разные цвета для разных функциональных групп
3. **Snap to variables** — при перетаскивании технологии привязывать к ближайшей переменной
4. **Предупреждения** — если координата не соответствует стандартной переменной

---

## Итоговая структура данных

```json
{
  "technology_id": "radio_technology",
  "folder": "electronics_folder",
  "position": {
    "x": -2,
    "y": 4,
    "x_var": "@HQ",
    "y_var": "@1937"
  },
  "start_year": 1937,
  "subtree": "Radio Communications",
  "research_cost": 1.5
}
```

---

## Справочные материалы

- Файл с технологиями: `common/technologies/electronic_mechanical_engineering.txt`
- Определения переменных всегда в начале блока `technologies = { ... }`
- Все координаты — целые числа или переменные, ссылающиеся на целые числа
