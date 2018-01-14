const path = require('path');
const webpack = require('webpack');

const FlowtypePlugin = require('flowtype-loader/plugin');

module.exports = {
  entry: './src/App.js',
  output: {
    path: path.resolve('..', 'dist'),
    filename: 'App.js',
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        loader: 'flowtype-loader',
        enforce: 'pre',
        exclude: /node_modules/,
      },
      {
        test: /\.(js|jsx)$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /\.(eot|woff|woff2|ttf|svg|png|jpg)$/,
        use: ['url-loader'],
      },
    ],
  },
  plugins: [
    new webpack.ProvidePlugin({
      jQuery: 'jquery',
    }),
    new FlowtypePlugin(),
  ],
};
