const path = require('path');

module.exports = {
  mode: 'development',
  devServer: { static: './dist' },
  devtool: 'inline-source-map',
  //optimization: { runtimeChunk: 'single' },
  entry: './src/index.js',
  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, 'dist'),
  },
  module: {
    rules: [
      { test: /\..json$/, use: { loader: 'json-loader' } },
      { test: /\.css$/i, use: [ 'style-loader','css-loader' ] }
    ]
  }
};
