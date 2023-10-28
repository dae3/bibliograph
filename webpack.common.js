const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  mode: 'development',
  devServer: { static: './dist' },
  devtool: 'inline-source-map',
  //optimization: { runtimeChunk: 'single' },
  entry: './src/index.js',
  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, 'dist'),
    clean: true
  },
  plugins: [
    new HtmlWebpackPlugin({ title: 'Bibliograph' }),
  ],
  module: {
    rules: [
      { test: /\..json$/, use: { loader: 'json-loader' } },
      { test: /\.css$/i, use: [ 'style-loader','css-loader' ] }
    ]
  }
};
