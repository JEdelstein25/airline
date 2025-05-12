import { workaroundForMultiTypeUnions } from '$lib/api'

// Disable SSR completely for now due to API connection issues
export const ssr = false

workaroundForMultiTypeUnions()
