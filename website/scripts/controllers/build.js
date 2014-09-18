(function () {

	/**
	 * Build controller responsible for displaying
	 * the build page and rendering the build results.
	 */	
	function BuildCtrl($scope, $routeParams, packages, builds) {
		var package = $routeParams.package;
		var version = $routeParams.version;

		// gets the build data from the server
		builds.get(package, version).then(function(build){
			$scope.build = build.data;
			$scope.version = version;
			var channel = build.data.channel
			var sdk     = build.data.sdk

			// gets the build output from the server
			builds.getOutput(package, version, channel, sdk).then(function(output){
				$scope.output = output.data;
			}).catch(function(error){
				$scope.error = error;
			});

		}).catch(function(error){
			$scope.error = error;
		});

		// gets the package data from the server
		packages.get(package).then(function(pkg){
			$scope.pkg = pkg.data;
		}).catch(function(error){
			$scope.error = error;
		});
	}

	angular
		.module('app')
		.controller('BuildCtrl', BuildCtrl);
})();