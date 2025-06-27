/**
 * 下载 DDNS-Go 到指定路径
 * 依赖：api.js, public.js (用于showNotification)
 */
async function downloadDDNSGo() {
  const DownLoadlink = document.getElementById('ThirdPartyExtDdnsGODownLoadlink').value
  const Version = document.getElementById('ThirdPartyExtDdnsGOVersion').value
  const DdnsGOBinPath = document.getElementById('DdnsGOBinPath').value
  const ThirdPartyExtGitHubDownloadMirror = document.getElementById('ThirdPartyExtGitHubDownloadMirror').value

  showNotification('DDNS-Go 正在下载，请不要离开页面', 'info')

  try {
    const response = await API.request(
      `/@api/admin/get_ThirdParty_ddnsgo?DownLoadlink=${encodeURIComponent(DownLoadlink)}&Version=${encodeURIComponent(Version)}&DdnsGOBinPath=${encodeURIComponent(DdnsGOBinPath)}&GitHubDownloadMirror=${encodeURIComponent(ThirdPartyExtGitHubDownloadMirror)}`,
      {},
      { needToken: true }
    )

    if (response.code < 10) {
      showNotification('DDNS-Go 下载成功', 'success')
    } else {
      showNotification('DDNS-Go 下载失败: ' + response.message, 'danger')
    }
  } catch (error) {
    showNotification('DDNS-Go 下载出错: ' + error.message, 'danger')
  }
}
