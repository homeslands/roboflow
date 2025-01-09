import antfu from '@antfu/eslint-config'

export default antfu(
  {
    vue: true,
    typescript: true,
  },
  {
    files: ['**/*.vue'],
    rules: {
      'vue/one-component-per-file': 0,
      'vue/no-reserved-component-names': 0,
      'vue/no-useless-v-bind': 0,
    },
  },
  {
    // Without `files`, they are general rules for all files
    rules: {
      'symbol-description': 0,
      'no-console': 'warn',
      'no-alert': 'warn',
      'no-tabs': 'warn',
      'ts/no-explicit-any': 'warn',
      'vue/no-unused-refs': 'warn',
      'unused-imports/no-unused-vars': 'warn',
      'import/first': 0,
      'node/prefer-global/process': 0,
      'style/no-tabs': 0,
      'unicorn/no-new-array': 0,
    },
  },
)
