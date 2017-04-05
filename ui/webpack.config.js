const path = require("path");
const webpack = require("webpack");

const FlowtypePlugin = require('flowtype-loader/plugin');

const config = {
  devtool: "#source-maps",
  entry: "./src/App.js",
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "App.js"
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        loader: 'flowtype-loader',
        enforce: 'pre',
        exclude: /node_modules/
      },
      {
        test: /\.(js|jsx)$/,
        loader: 'babel-loader',
        options: {
          presets: ["es2015", "react"],
          plugins: ["transform-flow-strip-types", "transform-class-properties"]
        },
        exclude: /node_modules/
      },
      {
        test: /\.css$/,
        use: ["style-loader", "css-loader"]
      },
      {
        test: /\.(eot|woff|woff2|ttf|svg|png|jpg)$/,
        use: ["url-loader"]
      }
    ]
  },
  plugins: [
    new webpack.ProvidePlugin({
      'jQuery': 'jquery'
    }),
    new FlowtypePlugin()
  ],
  devServer: {
    contentBase: path.join(__dirname, "dist")
  },
}

module.exports = config;
