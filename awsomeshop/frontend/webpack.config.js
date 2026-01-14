const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = (env, argv) => {
  const isProduction = argv.mode === 'production';

  return {
    // Entry point
    entry: './src/js/app.js',

    // Output configuration
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: isProduction ? 'js/bundle.[contenthash].js' : 'js/bundle.js',
      clean: true, // Clean dist folder before each build
      publicPath: '/',
    },

    // Module rules for different file types
    module: {
      rules: [
        // JavaScript files - Babel transpilation
        {
          test: /\.js$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            // Babel configuration is in .babelrc file
          },
        },
        // CSS files
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader'],
        },
        // Image files
        {
          test: /\.(png|jpg|jpeg|gif|svg)$/i,
          type: 'asset/resource',
          generator: {
            filename: 'images/[name].[hash][ext]',
          },
        },
      ],
    },

    // Plugins
    plugins: [
      new HtmlWebpackPlugin({
        template: './src/index.html',
        filename: 'index.html',
        inject: 'body',
      }),
    ],

    // Development server configuration
    devServer: {
      static: {
        directory: path.join(__dirname, 'dist'),
      },
      compress: true,
      port: 8000,
      hot: true,
      historyApiFallback: true, // Support History API routing
      open: true,
      proxy: [
        {
          context: ['/api'],
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      ],
    },

    // Source maps for debugging
    devtool: isProduction ? 'source-map' : 'eval-source-map',

    // Optimization
    optimization: {
      minimize: isProduction,
      splitChunks: isProduction ? {
        chunks: 'all',
        cacheGroups: {
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendors',
            priority: 10,
          },
          common: {
            minChunks: 2,
            priority: 5,
            reuseExistingChunk: true,
          },
        },
      } : false,
      runtimeChunk: isProduction ? 'single' : false,
    },

    // Resolve configuration
    resolve: {
      extensions: ['.js', '.json'],
      alias: {
        '@': path.resolve(__dirname, 'src'),
      },
    },

    // Performance hints
    performance: {
      hints: isProduction ? 'warning' : false,
      maxEntrypointSize: 512000,
      maxAssetSize: 512000,
    },
  };
};
