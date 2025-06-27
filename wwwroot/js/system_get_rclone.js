/**
 * 下载 rclone 到指定路径
 * 依赖：api.js, public.js (用于showNotification)
 */
async function downloadRclone() {
  const DownLoadlink = document.getElementById('ThirdPartyExtRcloneDownLoadlink').value
  const Version = document.getElementById('ThirdPartyExtRcloneVersion').value
  const BinPath = document.getElementById('RcloneBinPath').value
  const ThirdPartyExtGitHubDownloadMirror = document.getElementById('ThirdPartyExtGitHubDownloadMirror').value

  showNotification('DDNS-Go 正在下载，请不要离开页面', 'info')

  try {
    const response = await API.request(
      `/@api/admin/get_ThirdParty_rclone?DownLoadlink=${encodeURIComponent(DownLoadlink)}&Version=${encodeURIComponent(Version)}&BinPath=${encodeURIComponent(BinPath)}&GitHubDownloadMirror=${encodeURIComponent(ThirdPartyExtGitHubDownloadMirror)}`,
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
