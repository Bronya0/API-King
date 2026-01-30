<!--
  - Copyright 2025 Bronya0 <tangssst@163.com>.
  - Author Github: https://github.com/Bronya0
  -
  - Licensed under the Apache License, Version 2.0 (the "License");
  - you may not use this file except in compliance with the License.
  - You may obtain a copy of the License at
  -
  -     https://www.apache.org/licenses/LICENSE-2.0
  -
  - Unless required by applicable law or agreed to in writing, software
  - distributed under the License is distributed on an "AS IS" BASIS,
  - WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  - See the License for the specific language governing permissions and
  - limitations under the License.
  -->

<template>
  <n-page-header style="padding: 4px;--wails-draggable:drag">
    <template #avatar>
      <n-avatar :src="icon"/>
    </template>
    <template #subtitle>
      {{ props.desc }}

<!--      <n-tooltip>-->
<!--        <template #trigger>-->
<!--          <n-tag :type=title_tag v-if="subtitle">{{ subtitle }}</n-tag>-->
<!--          <n-p v-else>{{ desc }}</n-p>-->
<!--        </template>-->
<!--        健康：{{ title_tag }}-->
<!--      </n-tooltip>-->
    </template>
    <template #title>
      <div style="font-weight: 800">{{ app_name }}🎉</div>
    </template>
    <template #extra>
      <n-flex align="center" justify="flex-end" style="--wails-draggable:no-drag" class="right-section">
<!--        <n-button quaternary :focusable="false" @click="openUrl(qq_url)">技术交流群</n-button>-->
        <!--        <n-button quaternary :focusable="false" @click="changeTheme" :render-icon="renderIcon(MoonOrSunnyOutline)"/>-->

<!--        <n-tooltip placement="bottom" trigger="hover">-->
<!--          <template #trigger>-->
<!--            <n-button quaternary @click="openUrl(update_url)"-->
<!--                      :render-icon="renderIcon(HouseTwotone)"/>-->
<!--          </template>-->
<!--          <span>主页</span>-->
<!--        </n-tooltip>-->

<!--        <n-tooltip placement="bottom" trigger="hover">-->
<!--          <template #trigger>-->
<!--            <n-button quaternary :focusable="false" :loading="update_loading" @click="checkForUpdates"-->
<!--                      :render-icon="renderIcon(SystemUpdateAltSharp)"/>-->
<!--          </template>-->
<!--          <span>检查版本：{{ version.tag_name }} {{ check_msg }}</span>-->
<!--        </n-tooltip>-->
        <span>{{ version }}</span>
        <n-button quaternary class="ope_btn" :focusable="false" @click="minimizeWindow"
                  :render-icon="renderIcon(RemoveOutlined)"/>
        <n-button quaternary class="ope_btn" :focusable="false" @click="resizeWindow"
                  :render-icon="renderIcon(MaxMinIcon)"/>
        <n-button quaternary class="ope_btn close-btn" :focusable="false" @click="closeWindow"
                  :render-icon="renderIcon(CloseFilled)"></n-button>

      </n-flex>
    </template>
  </n-page-header>
</template>

<script setup>
import {NAvatar, NButton, NFlex} from 'naive-ui'
import icon from '../assets/images/appicon.png'
import {onMounted, ref, shallowRef} from "vue";
import {Window} from "@wailsio/runtime";
import {renderIcon} from "../utils/common";
import {GetAppName, GetVersion} from "../../bindings/app/backend/service/app";
import CropSquareFilled from "../assets/icons/CropSquareFilled.svg";
import ContentCopyFilled from "../assets/icons/ContentCopyFilled.svg";
import RemoveOutlined from "../assets/icons/RemoveOutlined.svg";
import CloseFilled from "../assets/icons/CloseFilled.svg";

const props = defineProps(['desc']);

const app_name = ref("");
const MaxMinIcon = shallowRef(CropSquareFilled)

const version = ref("")

onMounted(async () => {

  app_name.value = await GetAppName()
  version.value = await GetVersion()


})

const minimizeWindow = async () => {
  await Window.Minimise()
}

const resizeWindow = async () => {
  if (await Window.IsMaximised()) {
    await Window.UnMaximise();
    MaxMinIcon.value = CropSquareFilled;
  } else {
    await Window.Maximise();
    MaxMinIcon.value = ContentCopyFilled;
  }

}

const closeWindow = async () => {
  await Window.Hide()
}
// const changeTheme = () => {
//   MoonOrSunnyOutline.value = MoonOrSunnyOutline.value === NightlightRoundFilled ? WbSunnyOutlined : NightlightRoundFilled;
//   theme = MoonOrSunnyOutline.value === NightlightRoundFilled ? darkTheme : lightTheme
//   emitter.emit('update_theme', theme)
// }
</script>

<style scoped>


.right-section .n-button {
  padding: 0 8px;
}

.close-btn:hover {
  background-color: #800020;
}
.ope_btn:hover {
  transition: none !important; /* 移除所有过渡效果 */
}
</style>