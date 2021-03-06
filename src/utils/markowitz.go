// Author: Felipe Bogaerts de Mattos
// email: felipebogaerts@gmail.com
// License: MIT

package utils

import (
	"math/rand"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func MarkowitzOptimizer(df dataframe.DataFrame) dataframe.DataFrame {
	// Getting mean values:
	dfAnualisedReturns := GetAnualisedReturns(df)
	dfCovariance := GetCovariance(df)

	// Deleting date column:
	df = df.Drop(0)

	numberOfPortfolios := 1000        // number of portfolios to sample
	assetNames := df.Names()          // data column names or asset names
	numberOfAssets := len(assetNames) // number of assets
	rand.Seed(75)                     // fixed seed, for reproducibility

	var portfolioReturnSlice []float64
	var portfolioVarianceSlice []float64

	// Generate portfolios:
	for i := 0; i < numberOfPortfolios; i++ {

		// Choose random asset:
		randomAssetIndexes := rand.Perm(numberOfAssets)

		// Generate random weights
		var weights []float64
		var weightSum float64 = 0
		for j := 0; j < numberOfAssets; j++ {
			randomNumber := rand.Float64()
			weights = append(weights, randomNumber)
			weightSum += randomNumber
		}
		for j := 0; j < len(weights); j++ {
			weights[j] = weights[j] / weightSum
		}

		// Calculating portfolio return and variance:
		portfolioReturn := 0.0
		portfolioVariance := 0.0

		for j := 0; j < numberOfAssets; j++ {
			anualisedReturnsSlice := dfAnualisedReturns.Col(assetNames[randomAssetIndexes[j]]).Float()
			weightedAnualisedReturns := 0.0
			for k := 0; k < len(anualisedReturnsSlice); k++ {
				weightedAnualisedReturns += weights[j] * anualisedReturnsSlice[k]
			}
			portfolioReturn += weightedAnualisedReturns

			for k := 0; k < numberOfAssets; k++ {
				weightedPortfolioVariance := 0.0
				covarianceCol := dfCovariance.Col(assetNames[k])
				for l := 0; l < covarianceCol.Len(); l++ {
					weightedPortfolioVariance += weights[j] * weights[k] * covarianceCol.Elem(k).Float()
				}
				portfolioVariance += weightedPortfolioVariance
			}
		}

		portfolioReturnSlice = append(portfolioReturnSlice, portfolioReturn)
		portfolioVarianceSlice = append(portfolioVarianceSlice, portfolioVariance)

	}

	portfolioReturnVarianceDf := dataframe.New(
		series.New(portfolioReturnSlice, series.Float, "Return"),
		series.New(portfolioVarianceSlice, series.Float, "Variance"),
	)

	return portfolioReturnVarianceDf
}
