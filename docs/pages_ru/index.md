---
title: GitOps CLI-утилита
permalink: /
layout: default
---

<div class="presentation" id="presentation">
    <div class="presentation__bg" id="presentation-bg"></div>
    <div class="page__container presentation__container">
        <div class="presentation__row">
            <div class="presentation__row-item" id="presentation-title">
                <div class="presentation__subtitle">Инструмент консистентной доставки</div>
                <h1 class="presentation__title">What you Git<br/> is what you get!<span title="Что ты git'ишь, то и видишь!">*</span></h1>
                <ul class="presentation__features">
                    <li>Git — единый источник истины.</li>
                    <li>Сборка. Деплой в Kubernetes. Постоянная синхронизация.</li>
                    <li>Open Source. Написано на Go.</li>
                </ul>
            </div>
            <div class="presentation__row-item presentation__row-item_scheme">
                {% include scheme_ru.md %}
            </div>
        </div>
    </div>
</div>

<div class="welcome">
    <div class="page__container">
        <div class="welcome__content">
            <h1 class="welcome__title">
                Это GitOps, но <span>по-другому</span>!
            </h1>
            <div class="welcome__subtitle">
                Git, будучи единым источником истины, позволяет добиться детерминированного и&nbsp;идемпотентного процесса доставки по&nbsp;всему пайплайну. 
                Возможно применение как из&nbsp;CI-системы, так&nbsp;и&nbsp;с&nbsp;оператором (фича в&nbsp;разработке и&nbsp;будет доступна в&nbsp;ближайшее время).
            </div>
        </div>
    </div>
</div>

<div class="page__container">
    <div class="intro">
        <div class="intro__image"></div>        
    </div>
</div>

<div class="features">
    <div class="page__container">
        <ul class="features__list">
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_lifecycle"></div>
                <div class="features__list-item-title">Консольная утилита</div>
                <div class="features__list-item-text">
                    werf — это не SAAS, а самодостаточная CLI-утилита с открытым кодом, запускаемая на стороне клиента. Её можно использовать как для <b>локальной разработки</b>, так и для <b>встраивания в любую существующую CI/CD-систему</b>, оперируя основными командами как составляющими пайплайна:
                    <ul>
                        <li><code>werf converge</code>;</li>
                        <li><code>werf dismiss</code>;</li>
                        <li><code>werf cleanup</code>.</li>
                    </ul>
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_easy"></div>
                <div class="features__list-item-title">Простая в использовании</div>
                <div class="features__list-item-text">
                    werf работает «из коробки» с минимальной конфигурацией. Вам не нужно быть DevOps/SRE-инженером, чтобы использовать werf. Доступно <a href="{{ site.baseurl }}/documentation/guides.html"><b>множество гайдов</b></a>, которые помогут быстро организовать деплой приложений в Kubernetes и для целей разработки, и для production.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_config"></div>
                <div class="features__list-item-title">Объединяет лучшее</div>
                <div class="features__list-item-text">
                    werf связывает привычные инструменты, превращая их в понятую, целостную, <b>интегрированную CI/CD-платформу</b>. werf делает хорошо контролируемым и удобным взаимодействие Git, Docker, вашего container registry и существующей CI-системы, Helm и Kubernetes.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_kubernetes"></div>
                <div class="features__list-item-title">Распределенная сборка</div>
                <div class="features__list-item-text">
                    В werf реализован продвинутый сборщик, среди возможностей которого — алгоритм распределенной сборки. Благодаря нему и его распределенному кэшированию <b>ваши пайплайны становятся по-настоящему быстрыми</b>.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_debug"></div>
                <div class="features__list-item-title">Встроенная очистка</div>
                <div class="features__list-item-text">
                    Продуманный алгоритм <b>очистки неиспользуемых Docker-образов</b> в werf основан на анализе Git-истории собираемых приложений.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_helm"></div>
                <div class="features__list-item-title">Расширенный Helm</div>
                <div class="features__list-item-text">
                    В werf встроен бинарник <code>helm</code>, который реализует процесс деплоя, совместимый с Helm, и расширяет его возможности. С ним не требуется отдельная установка <code>helm</code>, а его дополнения обеспечивают детальные и понятные <b>логи при деплое</b>, быстрое <b>определение сбоев</b> во время деплоя, поддержку секретов и другие фичи, превращающие деплой в <b>понятный и надежный процесс</b>.
                </div>
            </li>
            <li class="features__list-item features__list-item_special">
                <div class="features__list-item-title">Open Source</div>
                <div class="features__list-item-description">
                    <a href="https://github.com/werf/werf" target="_blank">Код открыт</a> и написан на Go. За годы развития проекта у него сформировалось большое сообщество пользователей.
                </div>
            </li>
        </ul>
    </div>
</div>

<div class="stats">
    <div class="page__container">
        <div class="stats__content">
            <div class="stats__title">Активная разработка</div>
            <ul class="stats__list">
                <li class="stats__list-item">
                    <div class="stats__list-item-num">4</div>
                    <div class="stats__list-item-title">релиза в неделю</div>
                    <div class="stats__list-item-subtitle">в среднем за прошлый год</div>
                </li>
                <li class="stats__list-item">
                    <div class="stats__list-item-num">2000+</div>
                    <div class="stats__list-item-title">инсталляций</div>
                    <div class="stats__list-item-subtitle">в больших и маленьких проектах</div>
                </li>
                <li class="stats__list-item">
                    <div class="stats__list-item-num gh_counter">1470</div>
                    <div class="stats__list-item-title">звезд на GitHub</div>
                    <div class="stats__list-item-subtitle">поддержите проект ;)</div>
                </li>
            </ul>
        </div>
    </div>
</div>

<div class="reliability">
    <div class="page__container">
        <div class="reliability__content">
            <div class="reliability__column">
                <div class="reliability__title">
                    werf — это зрелый, надежный<br>
                    инструмент, которому можно доверять
                </div>
                <a href="{{ site.baseurl }}/installation.html#all-changes-in-werf-go-through-all-stability-channels" class="page__btn page__btn_b page__btn_small page__btn_inline">
                    Подробнее об уровнях стабильности и релизах
                </a>
            </div>
            <div class="reliability__column reliability__column_image">
                <div class="reliability__image"></div>
            </div>
        </div>
    </div>
</div>

<div class="community">
    <div class="page__container">
        <div class="community__content">
            <div class="community__title">Растущее дружелюбное сообщество</div>
            <div class="community__subtitle">Мы всегда на связи с сообществом<br/> в Telegram, Twitter и Discourse.</div>
            <div class="community__btns">
                <a href="{{ site.social_links[page.lang].telegram }}" target="_blank" class="page__btn page__btn_w community__btn">
                    <span class="page__icon page__icon_telegram"></span>
                    Мы в Telegram
                </a>
                <a href="{{ site.social_links[page.lang].twitter }}" target="_blank" class="page__btn page__btn_w community__btn">
                    <span class="page__icon page__icon_twitter"></span>
                    Мы в Twitter
                </a>
                <a href="https://community.flant.com/c/werf/6" rel="noopener noreferrer" target="_blank" class="page__btn page__btn_w community__btn">
                    <span class="page__icon page__icon_discourse"></span>
                    Мы в Discourse
                </a>
            </div>
        </div>
    </div>
</div>

<div class="page__container">
    <div class="documentation">
        <div class="documentation__image">
        </div>
        <div class="documentation__info">
            <div class="documentation__info-title">
                Исчерпывающая документация
            </div>
            <div class="documentation__info-text">
              <a href="{{ site.baseurl }}/documentation/index.html">Документация</a> содержит более 100 статей, включающих описание частых случаев (первые шаги, деплой в Kubernetes, интеграция с CI/CD-системами и другое), полное описание функций, архитектуры и CLI-команд.
            </div>
        </div>
        <div class="documentation__btns">
            <a href="{{ site.baseurl }}/introduction.html" class="page__btn page__btn_b documentation__btn">
                Знакомство
            </a>
            <a href="{{ site.baseurl }}/documentation/quickstart.html" class="page__btn page__btn_o documentation__btn">
                Начало работы
            </a>
            <a href="{{ site.baseurl }}/applications_guide_ru/" class="page__btn page__btn_o documentation__btn">
                Самоучители
            </a>
        </div>
    </div>
</div>
