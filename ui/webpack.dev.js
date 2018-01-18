const path = require('path');
const merge = require('webpack-merge');
const common = require('./webpack.common.js');
const env = require('./env.js');

const HtmlWebpackPlugin = require('html-webpack-plugin');
const InterpolateHtmlPlugin = require('interpolate-html-plugin');

module.exports = merge(common, {
  devtool: 'cheap-module-source-map',
  entry: './src/index.js',
  plugins: [
    new InterpolateHtmlPlugin(env),
    new HtmlWebpackPlugin({
      inject: false,
      template: './dist/index.html',
    }),
  ],
  devServer: {
    contentBase: path.join(__dirname, 'dist'),
    historyApiFallback: true,
  },
});
