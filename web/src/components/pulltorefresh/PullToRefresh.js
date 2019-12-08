import PullToRefresh from 'pulltorefreshjs'

const styles = () => {
    return `
        .__PREFIX__ptr {
            pointer-events: none;
            font-size: 0.85em;
            font-weight: bold;
            top: 0;
            height: 0;
            transition: height 0.3s, min-height 0.3s;
            text-align: center;
            width: 100%;
            overflow: hidden;
            display: flex;
            align-items: flex-end;
            align-content: stretch;
            font-family: -apple-system, BlinkMacSystemFont, helvetica neue, Helvetica, Arial, sans-serif;
            text-rendering: optimizeLegibility;
        }
        .__PREFIX__box {
            padding: 10px;
            flex-basis: 100%;
        }
        .__PREFIX__pull {
            transition: none;
        }
        .__PREFIX__text {
            margin-top: .33em;
            color: rgba(0, 0, 0, 0.3);
        }
        .__PREFIX__icon {
            color: rgba(0, 0, 0, 0.3);
            transition: transform .3s;
        }
        /*
        When at the top of the page, disable vertical overscroll so passive touch
        listeners can take over.
        */
        .__PREFIX__top {
            touch-action: pan-x pan-down pinch-zoom;
        }
        .__PREFIX__release .__PREFIX__icon {
            transform: rotate(180deg);
        }
    `
}

export default () => {
    const $ = {}

    $.onRefresh = () => {}

    $.shouldPullToRefresh = () => {
        return !window.scrollY
    }

    $.mainElement = 'main'

    $.init = () => {
        PullToRefresh.init({
            mainElement: $.mainElement,
            onRefresh() {
                $.onRefresh()
            },
            iconArrow: `
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-arrow-down">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <polyline points="19 12 12 19 5 12"/>
                </svg>
            `,
            iconRefreshing: `
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-loader">
                    <line x1="12" y1="2" x2="12" y2="6"/>
                    <line x1="12" y1="18" x2="12" y2="22"/>
                    <line x1="4.93" y1="4.93" x2="7.76" y2="7.76"/>
                    <line x1="16.24" y1="16.24" x2="19.07" y2="19.07"/>
                    <line x1="2" y1="12" x2="6" y2="12"/>
                    <line x1="18" y1="12" x2="22" y2="12"/>
                    <line x1="4.93" y1="19.07" x2="7.76" y2="16.24"/>
                    <line x1="16.24" y1="7.76" x2="19.07" y2="4.93"/>
                </svg>
            `,
            getStyles: styles,
            shouldPullToRefresh: $.shouldPullToRefresh,
        })
    }

    $.destroy = () => PullToRefresh.destroyAll()

    return $
}
