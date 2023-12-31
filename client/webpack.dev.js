const { merge } = require('webpack-merge');
const webpack  = require('webpack');
const common = require('./webpack.common.js');

module.exports = merge(common, {
  mode: 'development',
  devServer: { static: './dist' },
  devtool: 'inline-source-map',
  plugins: [ new webpack.DefinePlugin({
    BASE_URL: JSON.stringify("http://localhost:5555"),
    API_BASE: JSON.stringify("/api/v1/"),
    AUTH_BASE: JSON.stringify("/auth/"),
  }) ]
});
