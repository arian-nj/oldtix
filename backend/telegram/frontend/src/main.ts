import './style.css'
import WebApp from '@twa-dev/sdk'

WebApp.ready()

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
	<div class="game-container">
	  <iframe src="/game/mini.html" frameborder="0"></iframe>
	</div>
`
function getInitdata(message: String) {
	console.log(`message from getTelegram Said: ${message}`)
	return WebApp.initData
}

(window as any).getInitdata = getInitdata

