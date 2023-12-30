const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  optimization: {
    runtimeChunk: 'single',
    splitChunks: { chunks: 'all' }
  },
  entry: './src/index.js',
  output: {
    filename: '[chunkhash].js',
    path: path.resolve(__dirname, 'dist'),
    clean: true,
  },
  plugins: [
    new HtmlWebpackPlugin({ title: 'Bibliograph' }),
  ],
  module: {
    rules: [
      { test: /\..json$/, use: { loader: 'json-loader' } },
      { test: /\.css$/i, use: [ 'style-loader','css-loader', {
        loader: 'postcss-loader',
        options: {
          postcssOptions: { plugins: ['postcss-preset-env','tailwindcss','autoprefixer'] }
        }
      } ] },
      { test: /bib.json$/, type: 'asset/resource' },
      { test: /\.m?js$/, exclude: /node_modules/, use: { loader: 'babel-loader', options: { presets:  ['@babel/preset-react']}}}
    ]
  }
};
