/// <reference types="vite/client" />

interface AppWindowConfig {
  readonly logoLinkUrl?: string
}

interface Window {
  readonly __APP_CONFIG__?: AppWindowConfig
}

interface ImportMetaEnv {
  readonly LOGO_LINK_URL: string
}
