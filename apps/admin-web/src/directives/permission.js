import { useAuthStore } from '../stores/auth'

export function setupPermission(app) {
  app.directive('permission', {
    mounted(el, binding) {
      const auth = useAuthStore()
      const code = binding.value
      if (code && !auth.hasPerm(code)) {
        el.parentNode && el.parentNode.removeChild(el)
      }
    },
  })
}
