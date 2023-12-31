const { merge } = require('webpack-merge');
const webpack  = require('webpack');
const common = require('./webpack.common.js');

module.exports = merge(common, {
    mode: 'production',
    plugins: [ new webpack.DefinePlugin({
        BASE_URL: JSON.stringify("https://bibliograph.fly.dev"),
        API_BASE: JSON.stringify("/api/v1/"),
        AUTH_BASE: JSON.stringify("/auth/"),
    }) ]
});
