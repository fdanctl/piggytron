package charts

import (
	"fmt"

	"golang.org/x/text/currency"
)

// TODO: locale config
// use parseInt instead of addCommas for numeric values
// console.log(parseInt(params.value || '0', 10).toLocaleString('pt-PT'));

func pieTooltipFormatter(curr currency.Unit) string {
	return fmt.Sprintf(`
		function (params) {
			const formattedValue = new Intl.NumberFormat('en-US', {
				minimumFractionDigits: 2,
				maximumFractionDigits: 2
			}).format(params.value);
			return '<div style="color: var(--grey-200);">' +
				params.marker +
    		    params.name + '<br/>' +
				'<strong>' + formattedValue + '</strong>' +
				'<span style="font-size: 10px"> %s</span> (' + params.percent + '%%)'+
			'</div>';
		}`, curr)
}

func pieTooltipFormatterWithSeries(curr currency.Unit) string {
	return fmt.Sprintf(`
		function (params) {
			const formattedValue = new Intl.NumberFormat('en-US', {
				minimumFractionDigits: 2,
				maximumFractionDigits: 2
			}).format(params.value);
			return '<div style="color: var(--grey-200);">' +
				params.marker +
    		    params.name + ' (' + params.seriesName + ')' + '<br/>' +
				'<strong>' + formattedValue + '</strong>' +
				'<span style="font-size: 10px"> %s</span> (' + params.percent + '%%)'+
			'</div>';
		}`, curr)
}

func lineTooltipFormatter(curr currency.Unit) string {
	return fmt.Sprintf(`
		function (params) {
			const formater = echarts.format;
			const date = echarts.time.format(
				params[0].value[0],
    			'{yyyy}-{MM}-{dd}',
    			true
			);
		
			let strings = ['<div class="mb-sm">', date, '</div>'];
			strings = strings.concat(params.map((e) => {
				const formattedValue = new Intl.NumberFormat('en-US', {
					minimumFractionDigits: 2,
					maximumFractionDigits: 2
				}).format(e.value[1]);
				return (
					'<div class="grid gap-md grid-cols-2">' +
						'<span>' + e.marker + e.seriesName + ' ' + '</span>' +
						'<span class="text-right">' +
							'<strong>' + formattedValue + '</strong>' + '<span style="font-size: 10px"> %s</span>' +
						'</span>' +
					'</div>'
				);
			}));

			return '<div style="color: var(--grey-200);">' + strings.join('') + '</div>';
		}`, curr)
}

func barTooltipFormatter(curr currency.Unit) string {
	return fmt.Sprintf(`
		function (params) {
			const top = params.seriesName ? params.seriesName + '<br/>' : '';
			const marker = params.seriesName ? params.marker : '';
			const formattedValue = new Intl.NumberFormat('en-US', {
				minimumFractionDigits: 2,
				maximumFractionDigits: 2
			}).format(params.value);
			return '<div style="color: var(--grey-200);">' +
				 top +
				 marker + params.name +
				 '<strong class="ml-md">' + formattedValue + '</strong><span style="font-size: 10px"> %s</span>' + '</div>';
		}`, curr)
}

func sankeyTooltipFormatter(curr currency.Unit) string {
	return fmt.Sprintf(`
		function (params) {
			const formattedValue = new Intl.NumberFormat('en-US', {
				minimumFractionDigits: 2,
				maximumFractionDigits: 2
			}).format(params.value);
			const formater = echarts.format;
			let left = '';
			if (params.data.name) {
				left = params.data.name
			} else {
				left = params.data.source + ' → ' + params.data.target + ' '
			}
			return '<div style="color: var(--grey-200);">' +
				left +
				'<strong class="ml-md">' + formattedValue + '</strong>' + '<span style="font-size: 10px"> %s</span>' +
				'</div>';
		}`, curr)
}
