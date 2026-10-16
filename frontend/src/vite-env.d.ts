/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_BACKEND_URL?: string;
  readonly VITE_API_BASE_URL?: string;
  readonly VITE_APP_NAME?: string;
  readonly VITE_APP_DESCRIPTION?: string;
  readonly VITE_APP_TIMEZONE?: string;
  readonly VITE_APP_DEFAULT_LOCALE?: string;
  readonly VITE_APP_CURRENCY?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
