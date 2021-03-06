---
title: Обзор
permalink: documentation/index.html
sidebar: documentation
---

Документация содержит более 100 различных статей, включая наиболее типичные примеры использования werf, подробное описание функций, архитектуры и параметров вызова.

Мы рекомендуем начинать знакомство с раздела [**Руководства**]({{ site.baseurl }}/documentation/guides/installation.html):

- [Установка]({{ site.baseurl }}/documentation/guides/installation.html) содержит зависимости и возможные варианты установки.
- [Первые шаги]({{ site.baseurl }}/documentation/guides/getting_started.html) помогает начать использовать werf с существующим Dockerfile. Вы можете легко запустить werf в вашем проекте прямо сейчас.
- [Деплой в Kubernetes]({{ site.baseurl }}/documentation/guides/deploy_into_kubernetes.html) — краткий пример развертывания приложения в кластере Kubernetes.
- [Интеграция с CI/CD-системами]({{ site.baseurl }}/documentation/guides/generic_ci_cd_integration.html) — общий подход к интеграции werf с любой CI/CD-системой.
- [Интеграция с GitLab CI]({{ site.baseurl }}/documentation/guides/gitlab_ci_cd_integration.html) расскажет всё об интеграции с GitLab CI: про сборку, публикацию, деплой и очистку Docker registry.
- [Интеграция с GitHub Actions]({{ site.baseurl }}/documentation/guides/github_ci_cd_integration.html) расскажет всё об интеграции с GitHub Actions: про сборку, публикацию, деплой и очистку образов.
- В разделе расширенной сборки рассказывается о нашем синтаксисе описания сборки образов. Синтаксис позволяет использовать werf сборщик, который учитывает особенности инкрементальной сборки и предоставляет дополнительные возможности (к примеру, описание сборочных инструкций Ansible-задачами). Рекомендуем начать знакомство с создания [первого приложения]({{ site.baseurl }}/documentation/guides/advanced_build/first_application.html).

Следующий раздел — [**Конфигурация**]({{ site.baseurl }}/documentation/configuration/introduction.html).

Для использования werf в вашем проекте, необходимо создать файл конфигурации `werf.yaml`, который может состоять из:

1. Описания метаинформации проекта, которая впоследствии будет использоваться в большинстве команд и влиять на конечный результат (к примеру, на кэши и формат имён Helm-релиза и namespace в Kubernetes). Пример такой метаинформации — имя проекта.
2. Описания образов для сборки.

В статье [**Общие сведения**]({{ site.baseurl }}/documentation/configuration/introduction.html) вы найдете информацию о:

* Структуре секций и их конфигурации
* Описанию конфигурации в нескольких файлах
* Этапах обработки конфигурации 
* Поддерживаемых функциях Go-шаблонов

В других статьях раздела [**Конфигурация**]({{ site.baseurl }}/documentation/configuration/introduction.html) дается детальная информация о директивах описания [Dockerfile-образа]({{ site.baseurl }}/documentation/configuration/dockerfile_image.html), [Stapel-образа]({{ site.baseurl }}/documentation/configuration/stapel_image/naming.html), [Stapel-артефакта]({{ site.baseurl }}/documentation/internals/building_images_with_stapel/artifact.html) и особенностях их использования.

Раздел [**Справочник**]({{ site.baseurl }}/documentation/internals/building_of_images/images_storage.html) посвящен описанию основных процессов werf:

* [Сборка]({{ site.baseurl }}/documentation/reference/build_process.html)
* [Деплой]({{ site.baseurl }}/documentation/advanced/helm/basics.html)
* [Очистка]({{ site.baseurl }}/documentation/reference/cleaning_process.html)

Каждая статья описывает определенный процесс, особенности и доступные опции.

Также в этот раздел включены статьи с описанием базовых примитивов и общих инструментов:

* [Стадии и образы]({{ site.baseurl }}/documentation/internals/building_of_images/images_storage.html)
* [Работа с Docker registries]({{ site.baseurl }}/documentation/reference/working_with_docker_registries.html)
* [Разработка и отладка]({{ site.baseurl }}/documentation/reference/development_and_debug/setup_minikube.html)
* [Toolbox]({{ site.baseurl }}/documentation/reference/toolbox/slug.html)

Раздел [**CLI Commands**]({{ site.baseurl }}/documentation/reference/cli/werf_build.html) содержит как базовые, необходимые для управления процессом CI/CD, так и служебные команды, обеспечивающие расширенные функциональные возможности.

Раздел [**Разработка**]({{ site.baseurl }}/documentation/development/stapel.html) содержит информацию, предназначенную для более глубокого понимания работы werf.
