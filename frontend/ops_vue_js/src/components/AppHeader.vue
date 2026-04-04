<script setup>
import { ref, computed } from "vue";
import { useRouter, useRoute, RouterLink } from "vue-router";
import { useI18n } from "vue-i18n";
import { useUserStore } from "@/stores/user";
import {
  IconMoon,
  IconSun,
  IconLogout,
  IconUser,
  IconSettings,
  IconMenu2,
  IconX,
} from "@tabler/icons-vue";

const { t, locale } = useI18n();
const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const isDark = ref(document.documentElement.classList.contains("dark"));
const mobileMenuOpen = ref(false);
const userDropdownOpen = ref(false);

function toggleTheme() {
  isDark.value = !isDark.value;
  document.documentElement.classList.toggle("dark", isDark.value);
  localStorage.setItem("theme", isDark.value ? "dark" : "light"); // 使用统一的'theme' key
}

function toggleLocale() {
  const newLocale = locale.value === "zh-CN" ? "en" : "zh-CN";
  locale.value = newLocale;
  localStorage.setItem('locale', newLocale);
}

function isActive(path) {
  return route.path === path;
}

function handleLogout() {
  userStore.logout();
  router.push("/login");
}

const activeClass = "bg-blue-50 text-blue-600 dark:bg-dk-card dark:text-blue-400";
const normalClass = "rounded-md px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dk-subtle dark:hover:bg-dk-card dark:hover:text-dk-text";

const navItems = computed(() => [
  { label: t("appname.home"), to: "/" },
  { label: t("appname.schedule"), to: "/schedule" },
  { label: t("appname.purchase"), to: "/purchase" },
  { label: t("appname.warehouse"), to: "/warehouse" },
]);
</script>

<template>
  <header
    class="border-b border-gray-200 bg-white dark:border-dk-muted dark:bg-dk-base print:hidden"
  >
    <div class="mx-auto flex h-14 max-w-6xl items-center px-4">
      <!-- Logo -->
      <RouterLink to="/" class="mr-6 flex items-center">
        <img src="/logo.svg" class="h-8 w-8 rounded-lg" alt="Operations" />
        <span class="ml-2 text-lg font-bold text-gray-800 dark:text-dk-text"
          >Operations</span
        >
      </RouterLink>

      <!-- Mobile toggle -->
      <button
        class="ml-auto rounded p-1.5 text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-dk-card dark:hover:text-dk-text md:ml-0 md:hidden"
        @click="mobileMenuOpen = !mobileMenuOpen"
      >
        <IconX v-if="mobileMenuOpen" :size="22" />
        <IconMenu2 v-else :size="22" />
      </button>

      <!-- Desktop Nav -->
      <nav class="hidden flex-1 items-center gap-1 md:flex">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          :class="[normalClass, isActive(item.to) && activeClass]"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <!-- Right actions -->
      <div class="ml-auto hidden items-center gap-1 md:flex">
        <!-- Language -->
        <button
          class="rounded-md p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-dk-card dark:hover:text-dk-text"
          @click="toggleLocale"
          :title="locale === 'zh-CN' ? 'English' : '中文'"
        >
          <span class="text-xs font-semibold uppercase">{{
            locale === "zh-CN" ? "EN" : "中"
          }}</span>
        </button>

        <!-- Theme -->
        <button
          class="rounded-md p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-dk-card dark:hover:text-dk-text"
          @click="toggleTheme"
        >
          <IconMoon v-if="!isDark" :size="20" />
          <IconSun v-else :size="20" />
        </button>

        <!-- User -->
        <div v-if="userStore.isLoggedIn" class="relative ml-1">
          <button
            class="flex items-center gap-1.5 rounded-md px-2 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-100 dark:text-dk-subtle dark:hover:bg-dk-card"
            @click="userDropdownOpen = !userDropdownOpen"
          >
            <img
              :src="userStore.avatarUrl"
              class="h-6 w-6 rounded-full object-cover"
              alt="avatar"
            />
            <span class="max-w-24 truncate">{{
              userStore.userInfo?userStore.userInfo.Username:userStore.user.Name
            }}</span>
          </button>
          <Transition
            enter-active-class="transition duration-100 ease-out"
            enter-from-class="transform scale-95 opacity-0"
            enter-to-class="transform scale-100 opacity-100"
            leave-active-class="transition duration-75 ease-in"
            leave-from-class="transform scale-100 opacity-100"
            leave-to-class="transform scale-95 opacity-0"
          >
            <div
              v-if="userDropdownOpen"
              class="absolute right-0 z-50 mt-1 w-48 rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              <RouterLink
                to="/settings/account"
                class="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-dk-subtle dark:hover:bg-dk-muted"
                @click="userDropdownOpen = false"
              >
                <IconSettings :size="16" />
                {{ t("message.user_settings") }}
              </RouterLink>
              <hr class="my-1 border-gray-200 dark:border-dk-muted" />
              <button
                class="flex w-full items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
                @click="
                  handleLogout();
                  userDropdownOpen = false;
                "
              >
                <IconLogout :size="16" />
                {{ t("message.logout") }}
              </button>
            </div>
          </Transition>
        </div>
        <RouterLink
          v-else
          to="/login"
          class="ml-2 rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          {{ t("message.login_or_register") }}
        </RouterLink>
      </div>

      <!-- Mobile user avatar (only when logged in) -->
      <div v-if="userStore.isLoggedIn" class="ml-3 md:hidden">
        <button
          class="rounded-md p-1.5 text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-dk-card dark:hover:text-dk-text"
          @click="userDropdownOpen = !userDropdownOpen"
        >
          <img
            :src="userStore.avatarUrl"
            class="h-7 w-7 rounded-full object-cover"
            alt="avatar"
          />
        </button>
        <Transition
          enter-active-class="transition duration-100 ease-out"
          enter-from-class="transform scale-95 opacity-0"
          enter-to-class="transform scale-100 opacity-100"
          leave-active-class="transition duration-75 ease-in"
          leave-from-class="transform scale-100 opacity-100"
          leave-to-class="transform scale-95 opacity-0"
        >
          <div
            v-if="userDropdownOpen"
            class="absolute right-4 z-50 mt-1 w-48 rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dk-muted dark:bg-dk-card"
          >
            <div class="px-4 py-2 text-sm text-gray-600 dark:text-dk-subtle">
              <div class="flex items-center gap-2">
                <img
                  :src="userStore.avatarUrl"
                  class="h-8 w-8 rounded-full object-cover"
                  alt="avatar"
                />
                <div class="truncate">{{ userStore.user?.Name || "" }}</div>
              </div>
            </div>
            <hr class="my-1 border-gray-200 dark:border-dk-muted" />
            <RouterLink
              to="/settings/account"
              class="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-dk-subtle dark:hover:bg-dk-muted"
              @click="userDropdownOpen = false"
            >
              <IconSettings :size="16" />
              {{ t("message.user_settings") }}
            </RouterLink>
            <hr class="my-1 border-gray-200 dark:border-dk-muted" />
            <button
              class="flex w-full items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
              @click="
                handleLogout();
                userDropdownOpen = false;
              "
            >
              <IconLogout :size="16" />
              {{ t("message.logout") }}
            </button>
          </div>
        </Transition>
      </div>
    </div>

    <!-- Mobile menu -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="-translate-y-2 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="-translate-y-2 opacity-0"
    >
      <div
        v-if="mobileMenuOpen"
        class="border-t border-gray-200 bg-white px-4 pb-4 pt-2 dark:border-dk-muted dark:bg-dk-base md:hidden"
      >
        <nav class="flex flex-col gap-1">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            :class="['rounded-md px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 dark:text-dk-subtle dark:hover:bg-dk-card', isActive(item.to) && 'bg-blue-50 text-blue-600 dark:bg-dk-card dark:text-blue-400']"
            @click="mobileMenuOpen = false"
          >
            {{ item.label }}
          </RouterLink>
        </nav>
        <hr class="my-3 border-gray-200 dark:border-dk-muted" />
        <div class="flex items-center gap-2">
          <button
            class="rounded-md p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-dk-card"
            @click="toggleLocale"
          >
            <span class="text-xs font-semibold uppercase">{{
              locale === "zh-CN" ? "EN" : "中"
            }}</span>
          </button>
          <button
            class="rounded-md p-2 text-gray-500 hover:bg-gray-100 dark:hover:bg-dk-card"
            @click="toggleTheme"
          >
            <IconMoon v-if="!isDark" :size="20" />
            <IconSun v-else :size="20" />
          </button>
<div class="ml-auto">
          <RouterLink
            v-if="!userStore.isLoggedIn"
            to="/login"
            class="rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700"
            @click="mobileMenuOpen = false"
          >
            {{ t("message.login_or_register") }}
          </RouterLink>
          <div v-else class="flex items-center gap-2 text-sm text-gray-600 dark:text-dk-subtle">
            <img
              :src="userStore.avatarUrl"
              class="h-6 w-6 rounded-full object-cover"
              alt="avatar"
            />
            <span class="truncate">{{ userStore.user?.Name || "" }}</span>
          </div>
        </div>
        </div>
      </div>
    </Transition>
  </header>
</template>
