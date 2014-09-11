'use script';

(function () {

	/**
	 * Creates the angular application.
	 */
	angular.module('app', [
			'ngRoute',
			'ui.filters'
		]);

	/**
	 * Defines the route configuration for the
	 * main application.
	 */
	function Config ($routeProvider) {
		$routeProvider
		.when('/:package/:version', {
			templateUrl: '/scripts/views/build.html',
			controller: 'BuildCtrl'
		})
		.when('/', {
			templateUrl: '/scripts/views/main.html',
			controller: 'MainCtrl'
		});
	}

	angular
		.module('app')
		.config(Config);

})();