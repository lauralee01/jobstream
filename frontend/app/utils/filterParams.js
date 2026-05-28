export const DEFAULT_SORT_BY = 'posted_at'
export const DEFAULT_SORT_ORDER = 'desc'
export const PAGE_SIZE = 20

export const emptyFilters = () => ({
  keyword: '',
  location: '',
  category: '',
  platforms: [],
  remote: false,
  salaryMin: '',
  sortBy: DEFAULT_SORT_BY,
  sortOrder: DEFAULT_SORT_ORDER,
  page: 1
})

export const filtersFromQuery = (query) => ({
  keyword: query.q || '',
  location: query.location || '',
  category: query.category || '',
  platforms: query.platforms ? String(query.platforms).split(',') : [],
  remote: query.work_model === 'remote',
  salaryMin: query.salary_min || query.min_salary || '',
  sortBy: query.sort_by || DEFAULT_SORT_BY,
  sortOrder: query.sort_order || DEFAULT_SORT_ORDER,
  page: parseInt(query.page, 10) || 1
})

export const queryFromFilters = (filters) => {
  const query = {}

  const keyword = filters.keyword?.trim()
  const location = filters.location?.trim()

  if (keyword) query.q = keyword
  if (location) query.location = location
  if (filters.category) query.category = filters.category
  if (filters.platforms?.length) query.platforms = filters.platforms.join(',')
  if (filters.remote) query.work_model = 'remote'
  if (filters.salaryMin) query.min_salary = String(filters.salaryMin)
  if (filters.page > 1) query.page = String(filters.page)
  if (filters.sortBy && filters.sortBy !== DEFAULT_SORT_BY) {
    query.sort_by = filters.sortBy
  }
  if (filters.sortOrder && filters.sortOrder !== DEFAULT_SORT_ORDER) {
    query.sort_order = filters.sortOrder
  }

  return query
}

/** Normalize filter object before applying (trim strings, reset page). */
export const normalizeFilters = (filters, { resetPage = true } = {}) => ({
  keyword: filters.keyword?.trim() || '',
  location: filters.location?.trim() || '',
  category: filters.category || '',
  platforms: filters.platforms || [],
  remote: Boolean(filters.remote),
  salaryMin: filters.salaryMin || '',
  sortBy: filters.sortBy || DEFAULT_SORT_BY,
  sortOrder: filters.sortOrder || DEFAULT_SORT_ORDER,
  page: resetPage ? 1 : (filters.page || 1)
})

export const buildJobsApiQuery = (filters) => ({
  keyword: filters.keyword || undefined,
  location: filters.location || undefined,
  category: filters.category || undefined,
  min_salary: filters.salaryMin || undefined,
  platforms: filters.platforms?.length ? filters.platforms.join(',') : undefined,
  remote: filters.remote ? 'true' : undefined,
  sort_by: filters.sortBy || DEFAULT_SORT_BY,
  sort_order: filters.sortOrder || DEFAULT_SORT_ORDER,
  page: filters.page || 1,
  limit: PAGE_SIZE
})
