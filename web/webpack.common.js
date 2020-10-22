const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = {
  plugins: [
    new MiniCssExtractPlugin({
      filename: '[name].[contenthash].css',
      chunkFilename: '[id].[contenthash].css'
    })
  ],
  entry: {
    // JS
    common: './static/js/common.js',
    signup: './static/js/signup.js',
    signin: './static/js/signin.js',
    // CSS
    style_default: './static/css/default.css',
  },
  output: {
    filename: '[name].[contenthash].js',
    path: path.resolve(__dirname, 'static/dist'),
  },

module: {
    rules: [
      {
        test: /\.m?js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env'],
            "plugins": [
                "@babel/plugin-proposal-class-properties"
              ]
          }
        }
      },
      {
        test: /\.css$/,
        use: [
          {
            loader: MiniCssExtractPlugin.loader,
            options: {
              publicPath: './static/dist/',
            },
          },
          'css-loader',
        ],
      },
      {
        test: /\.(woff(2)?|ttf|eot|svg|jpg|png)(\?v=\d+\.\d+\.\d+)?$/,
        use: [
          {
            loader: 'file-loader',
            options: {
              name: '[name].[ext]',
              outputPath: '.'
            }
          }
        ]
      }
    ]
  }
};
