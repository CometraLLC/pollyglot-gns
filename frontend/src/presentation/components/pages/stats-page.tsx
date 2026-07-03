'use client'

import { useQuery } from '@tanstack/react-query'
import { Bar, BarChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import { Flame, CalendarCheck, Library, BookMarked } from 'lucide-react'
import { statsService } from '@/src/domain/services/stats.service'
import type { LucideIcon } from 'lucide-react'

function StatTile({
	icon: Icon,
	label,
	value,
	suffix,
}: {
	icon: LucideIcon
	label: string
	value: number
	suffix?: string
}) {
	return (
		<div className='neu-card p-6' aria-label={`${label}: ${value}`}>
			<Icon className='mb-3 h-5 w-5 text-emerald-600 dark:text-emerald-400' aria-hidden />
			<p className='text-2xl font-bold tabular-nums'>
				{value}
				{suffix && <span className='ml-1 text-sm font-normal text-muted-foreground'>{suffix}</span>}
			</p>
			<p className='text-sm text-muted-foreground'>{label}</p>
		</div>
	)
}

export function StatsPage() {
	const { data, isPending, isError } = useQuery({
		queryKey: ['stats'],
		queryFn: () => statsService.getStats(),
	})

	return (
		<div className='mx-auto max-w-4xl'>
			<div className='mb-8'>
				<h1 className='text-2xl font-bold tracking-tight'>Progress</h1>
				<p className='text-muted-foreground'>Streaks, reviews, and every word you have met.</p>
			</div>

			{isPending && <p className='text-muted-foreground'>Loading your progress…</p>}

			{isError && (
				<p className='text-red-600 dark:text-red-400'>
					Could not load your progress. Check that the API is running, then reload.
				</p>
			)}

			{data && (
				<>
					<div className='mb-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
						<StatTile icon={Flame} label='Day streak' value={data.streak_days} />
						<StatTile icon={CalendarCheck} label='Reviews today' value={data.reviews_today} />
						<StatTile icon={Library} label='Total reviews' value={data.total_reviews} />
						<StatTile icon={BookMarked} label='Unique words' value={data.unique_cards} />
					</div>

					{data.total_reviews === 0 ? (
						<div className='neu-inset rounded-xl p-12 text-center'>
							<p className='mb-1 font-medium'>No reviews yet</p>
							<p className='text-sm text-muted-foreground'>
								Study a deck and this page starts filling in.
							</p>
						</div>
					) : (
						<div className='neu-card p-6'>
							<h2 className='mb-4 text-sm font-semibold'>Reviews per day — last 30 days</h2>
							{/* Chart color validated for both surfaces (dataviz six checks):
							    #059669 on light, #0ea371 on dark, inherited via currentColor. */}
							<div className='h-48 text-[#059669] dark:text-[#0ea371]'>
								<ResponsiveContainer width='100%' height='100%'>
									<BarChart data={data.reviews_per_day} barCategoryGap={2}>
										<XAxis
											dataKey='date'
											tickFormatter={(date: string) => date.slice(5)}
											interval={6}
											tickLine={false}
											axisLine={false}
											tick={{ fontSize: 11, fill: 'var(--muted-foreground)' }}
										/>
										<YAxis
											allowDecimals={false}
											width={28}
											tickLine={false}
											axisLine={false}
											tick={{ fontSize: 11, fill: 'var(--muted-foreground)' }}
										/>
										<Tooltip
											cursor={{ fill: 'transparent' }}
											contentStyle={{
												background: 'var(--card)',
												border: '1px solid var(--border)',
												borderRadius: 8,
												color: 'var(--foreground)',
												fontSize: 12,
											}}
											formatter={(value) => [`${value} reviews`, undefined]}
										/>
										<Bar
											dataKey='count'
											fill='currentColor'
											radius={[4, 4, 0, 0]}
											maxBarSize={16}
										/>
									</BarChart>
								</ResponsiveContainer>
							</div>
							{/* Accessible data table for the same series */}
							<table className='sr-only' aria-label='Reviews per day'>
								<thead>
									<tr>
										<th>Date</th>
										<th>Reviews</th>
									</tr>
								</thead>
								<tbody>
									{data.reviews_per_day.map((day) => (
										<tr key={day.date}>
											<td>{day.date}</td>
											<td>{day.count}</td>
										</tr>
									))}
								</tbody>
							</table>
						</div>
					)}
				</>
			)}
		</div>
	)
}
