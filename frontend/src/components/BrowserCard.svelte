<script>
  export let instance;
  export let minimalMode = false;
  export let proxies = [];
  export let onStart;
  export let onStop;
  export let onEdit;
  export let onDuplicate;
  export let onDelete;

  let selectedProxyId = 'none';
  let _lastInstId = null;
  $: if (instance && instance.id !== _lastInstId) {
    selectedProxyId = instance.proxyId || 'none';
    _lastInstId = instance.id;
  }

  $: proxyName = (() => {
    if (!instance.proxyId) return '';
    const p = proxies.find(p => p.id === instance.proxyId);
    return p ? p.name : '未知代理';
  })();

  $: statusColor = instance.running ? 'var(--success-color)' : 'var(--text-muted)';
  $: statusText = instance.running ? '运行中' : '已停止';
</script>

<div class="card" class:mini={minimalMode}>
  <div class="card-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
    <div style="flex: 1; min-width: 0; margin-right: 0.5rem;">
      <h3 style="font-size: 1.125rem; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;" title="{instance.sortNum}.{instance.name}">{instance.sortNum}.{instance.name}</h3>
      <div style="display: flex; align-items: center; gap: 0.5rem; margin-top: 0.25rem;">
        <span style="width: 8px; height: 8px; border-radius: 50%; background-color: {statusColor}"></span>
        <span style="font-size: 0.75rem; color: var(--text-muted); font-weight: 500;">{statusText}</span>
      </div>
    </div>
    

    
    <div class="actions" style="display: flex; gap: 0.5rem;">
      <button class="btn-icon" on:click={() => onEdit(instance)} title="编辑">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
      </button>
      <button class="btn-icon" on:click={() => onDuplicate(instance)} title="复制配置">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
      </button>
      <button class="btn-icon" on:click={() => onDelete(instance.id)} title="删除" style="color: var(--danger-color);">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
      </button>
    </div>
  </div>

  <div class="card-body" style="margin-bottom: 1.5rem;">
    {#if !minimalMode}
      <div class="info-item">
        <span class="label">启动路径:</span>
        <span class="value" title={instance.path}>{instance.path}</span>
      </div>
      <div class="info-item">
        <span class="label">数据目录:</span>
        <span class="value" title={instance.userDataDir}>{instance.userDataDir || '默认'}</span>
      </div>
      <div class="info-item">
        <span class="label">参数:</span>
        <span class="value">{instance.args.join(' ') || '无'}</span>
      </div>
      <div class="info-item">
        <span class="label">代理:</span>
        <span class="value">{proxyName || '不使用代理'}</span>
      </div>
    {/if}

      <div class="info-item" style="min-height: 2.5em;">
        <span class="label">标签:</span>
        <div class="tags-container">
          {#if instance.tags && instance.tags.length > 0}
            {#each instance.tags as tag}
               <span class="tag-badge">{tag}</span>
            {/each}
          {:else}
            <span class="text-muted" style="font-size: 0.75rem;">-</span>
          {/if}
        </div>
      </div>

  </div>

  <div class="card-footer">
    {#if !instance.running}
      <div class="start-group">
        <select class="proxy-select" bind:value={selectedProxyId}>
          <option value="none">不使用代理</option>
          {#each proxies as p}
            <option value={p.id}>{p.name}</option>
          {/each}
        </select>
        <button class="btn btn-primary" style="flex: 1;" on:click={() => onStart(instance.id, selectedProxyId)}>
          启动
        </button>
      </div>
    {:else}
      <button class="btn" style="flex: 1; border-color: var(--danger-color); color: var(--danger-color);" on:click={() => onStop(instance.id)}>
        停止
      </button>
    {/if}
  </div>

</div>

<style>
  .btn-icon {
    background: transparent;
    border: none;
    padding: 0.25rem;
    cursor: pointer;
    color: var(--text-muted);
    border-radius: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .btn-icon:hover {
    background: #f1f5f9;
    color: var(--text-main);
  }
  .info-item {
    display: flex;
    flex-direction: column;
    margin-bottom: 0.5rem;
  }
  .label {
    font-size: 0.75rem;
    color: var(--text-muted);
    margin-bottom: 0.125rem;
  }
  .value {
    font-size: 0.8125rem;
    color: var(--text-main);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }
  .tag-badge {
    background-color: #f1f5f9;
    color: var(--text-muted);
    padding: 0.125rem 0.375rem;
    border-radius: 0.25rem;
    font-size: 0.625rem;
  }
  .start-group {
    display: flex;
    gap: 0.5rem;
    width: 100%;
  }
  .proxy-select {
    padding: 0.375rem 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 0.375rem;
    font-size: 0.75rem;
    color: var(--text-main);
    background: white;
    cursor: pointer;
    min-width: 0;
    flex: 1;
    appearance: auto;
  }

  /* Minimal Mode Styles */
  .card.mini .card-header {
    margin-bottom: 0 !important; /* Remove bottom margin in header */
    padding-bottom: 0.5rem; /* Reduce padding inside header if needed */
  }
  .card.mini .card-body {
    margin-bottom: 0.5rem !important; /* Reduce body bottom margin */
  }
  .card.mini h3 {
    font-size: 1rem !important; /* Smaller title */
  }
  .card.mini .info-item {
    margin-bottom: 0.25rem !important; /* Compact items */
  }
</style>
