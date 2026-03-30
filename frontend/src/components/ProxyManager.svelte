<script>
  import { createEventDispatcher } from 'svelte';
  import ProxyModal from './ProxyModal.svelte';
  
  const dispatch = createEventDispatcher();

  export let show = false;
  export let proxies = [];
  export let onCheckProxy = null;

  let showProxyModal = false;
  let editingProxy = null;
  let checkResults = {}; // { [proxyId]: { region, latency, detail, loading } }

  function openAddProxy() {
    editingProxy = null;
    showProxyModal = true;
  }

  function openEditProxy(proxy) {
    editingProxy = JSON.parse(JSON.stringify(proxy));
    showProxyModal = true;
  }

  function handleProxySave(event) {
    dispatch('save', event.detail);
    showProxyModal = false;
  }

  function handleDeleteProxy(id) {
    if (!confirm("确定要删除此代理配置吗？")) return;
    dispatch('delete', { id });
  }

  function handleClose() {
    dispatch('close');
  }

  function formatProxyURL(proxy) {
    let url = proxy.protocol + '://';
    if (proxy.username) {
      url += proxy.username + ':***@';
    }
    url += proxy.host + ':' + proxy.port;
    return url;
  }

  async function handleCheckProxy(proxyId, target) {
    checkResults[proxyId] = { ...(checkResults[proxyId] || {}), loading: true };
    checkResults = checkResults; // trigger reactivity
    try {
      if (onCheckProxy) {
        const result = await onCheckProxy(proxyId, target);
        checkResults[proxyId] = { ...result, loading: false };
        checkResults = checkResults;
      }
    } catch (err) {
      checkResults[proxyId] = { error: err.toString(), loading: false };
      checkResults = checkResults;
    }
  }

  function getLatencyColor(latency) {
    if (!latency) return 'var(--text-muted)';
    if (latency < 100) return 'var(--success-color)';
    if (latency < 300) return '#f59e0b';
    return 'var(--danger-color)';
  }
</script>

{#if show}
  <div class="modal-overlay" on:click|self={handleClose}>
    <div class="modal" style="max-width: 650px;">
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem;">
        <h2 style="font-size: 1.25rem;">代理管理</h2>
        <button class="btn btn-primary btn-sm" on:click={openAddProxy}>
          <svg style="margin-right: 0.35rem;" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
          添加代理
        </button>
      </div>
      
      {#if proxies.length === 0}
        <div class="empty-state">
          还没有代理配置，点击"添加代理"创建一个吧。
        </div>
      {:else}
        <div class="proxy-list">
          {#each proxies as proxy (proxy.id)}
            <div class="proxy-item">
              <div class="proxy-info">
                <div class="proxy-name">{proxy.name}</div>
                <div class="proxy-url">{formatProxyURL(proxy)}</div>
                {#if checkResults[proxy.id]}
                  <div class="proxy-check-result">
                    {#if checkResults[proxy.id].loading}
                      <span class="checking">检测中...</span>
                    {:else if checkResults[proxy.id].error}
                      <span class="check-error">失败: {checkResults[proxy.id].error}</span>
                    {:else}
                      <span class="check-region">{checkResults[proxy.id].region}</span>
                      <span class="check-latency" style="color: {getLatencyColor(checkResults[proxy.id].latency)}">{checkResults[proxy.id].latency}ms</span>
                      {#if checkResults[proxy.id].detail}
                        <span class="check-detail" title={checkResults[proxy.id].detail}>ℹ️</span>
                      {/if}
                    {/if}
                  </div>
                {/if}
              </div>
              <div class="proxy-actions">
                <button class="btn-xs" title="检测 CN" on:click={() => handleCheckProxy(proxy.id, 'cn')}
                  disabled={checkResults[proxy.id]?.loading}>CN</button>
                <button class="btn-xs" title="检测 Global" on:click={() => handleCheckProxy(proxy.id, 'global')}
                  disabled={checkResults[proxy.id]?.loading}>GL</button>
                <button class="icon-btn" title="编辑" on:click={() => openEditProxy(proxy)}>
                  <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                </button>
                <button class="icon-btn danger" title="删除" on:click={() => handleDeleteProxy(proxy.id)}>
                  <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                </button>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      <div style="display: flex; justify-content: flex-end; margin-top: 1.5rem;">
        <button class="btn" on:click={handleClose}>关闭</button>
      </div>
    </div>
  </div>
{/if}

<ProxyModal
  show={showProxyModal}
  proxy={editingProxy}
  on:save={handleProxySave}
  on:close={() => showProxyModal = false}
/>

<style>
  .proxy-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    max-height: 400px;
    overflow-y: auto;
  }
  .proxy-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    background: #f8fafc;
    border: 1px solid var(--border-color);
    border-radius: 0.5rem;
    transition: background 0.15s;
  }
  .proxy-item:hover {
    background: #f1f5f9;
  }
  .proxy-info {
    flex: 1;
    min-width: 0;
  }
  .proxy-name {
    font-weight: 600;
    font-size: 0.9rem;
    color: var(--text-main);
  }
  .proxy-url {
    font-size: 0.75rem;
    color: var(--text-muted);
    font-family: monospace;
    margin-top: 0.125rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .proxy-check-result {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-top: 0.25rem;
    font-size: 0.75rem;
  }
  .checking {
    color: var(--primary-color);
    font-style: italic;
  }
  .check-error {
    color: var(--danger-color);
  }
  .check-region {
    font-weight: 600;
    color: var(--text-main);
  }
  .check-latency {
    font-weight: 500;
  }
  .check-detail {
    cursor: help;
    font-size: 0.7rem;
  }
  .proxy-actions {
    display: flex;
    gap: 0.375rem;
    flex-shrink: 0;
    margin-left: 0.75rem;
    align-items: center;
  }
  .icon-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    color: #64748b;
    padding: 0.375rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 0.25rem;
    transition: all 0.2s;
  }
  .icon-btn:hover {
    background-color: #e2e8f0;
    color: #334155;
  }
  .icon-btn.danger:hover {
    background-color: #fee2e2;
    color: #ef4444;
  }
  .btn-xs {
    padding: 0.125rem 0.375rem;
    font-size: 0.7rem;
    background-color: white;
    border: 1px solid var(--border-color);
    border-radius: 0.25rem;
    cursor: pointer;
    color: var(--primary-color);
    font-weight: 500;
  }
  .btn-xs:hover {
    background-color: #f8fafc;
  }
  .btn-xs:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .empty-state {
    text-align: center;
    padding: 2.5rem 1rem;
    color: var(--text-muted);
    background: #f8fafc;
    border: 1px dashed var(--border-color);
    border-radius: 0.5rem;
    font-size: 0.875rem;
  }
  .btn-sm {
    padding: 0.375rem 0.75rem;
    font-size: 0.8rem;
  }
</style>
