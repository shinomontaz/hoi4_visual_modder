Механизм связывания картинок с национальными фокусами и технологиями в Hearts of Iron IV
Национальные фокусы
Каждый фокус описан в файлах:
common/national_focus/*.txt

Картинка фокуса назначается через поле:

text
icon = "GFX_focus_my_icon"
Иконка связывается с графическим ресурсом, определённым в:

text
interface/goals.gfx
Пример spriteType:

text
spriteType = {
  name = GFX_focus_my_icon
  texturefile = "gfx/interface/goals/my_icon.dds"
}
Все файлы картинок для фокусов находятся по пути:
gfx/interface/goals/

Можно сделать динамические иконки через блок иконок с триггерами:

text
icon = {
  trigger = { <условие> }
  value = GFX_focus_icon_special
}
icon = { value = GFX_focus_icon_default }
dynamic = yes
Технологии
Описание технологий — файлы:
common/technologies/*.txt

Отдельного поля icon нет.

Для назначения картинки используется механизм spriteType в одном из файлов:

text
interface/countrytechtreeview.gfx
Пример записи:

text
spriteType = {
  name = GFX_technology_name_medium
  texturefile = "gfx/interface/technologies/my_tech_icon.dds"
}
Картинка для технологии ищется по имени:

Имя спрайта: GFX_<technology_id>_medium

Файл: gfx/interface/technologies/<имя файла>.dds

Для страновых картинок между GFX и именем технологии вставляется тег страны:

text
name = GFX_GER_infantry_equipment_2_medium
Тогда картинка применяется только для выбранной страны.

Структура связи
Элемент	Файл описания	Ключ для иконки	Файл спрайта	Каталог картинок
Национальный фокус	national_focus/*.txt	icon	interface/goals.gfx	gfx/interface/goals/
Технология	technologies/*.txt	нет отдельного поля	interface/countrytechtreeview.gfx	gfx/interface/technologies/
Итоги
Для национальных фокусов задаём картинку напрямую через поле icon и привязку в goals.gfx.

Для технологий картинка задаётся спрайтом в countrytechtreeview.gfx, а имя технологии используется для поиска нужного спрайта и картинки.

Картинки должны лежать в правильных каталогах:
gfx/interface/goals/ — фокусы
gfx/interface/technologies/ — технологии