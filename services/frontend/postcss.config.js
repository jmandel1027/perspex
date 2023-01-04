module.exports = {
  plugins: {
    tailwindcss: {},
    'tailwindcss/nesting': {},
    autoprefixer: {},
    'postcss-import': {},
    'postcss-preset-env': {
        features: { 'nesting-rules': false },
      },
    ...(process.env.NODE_ENV === 'production' ? { cssnano: {} } : {})
  }
}
