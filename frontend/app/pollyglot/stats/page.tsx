import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { StatsPage } from '@/src/presentation/components/pages/stats-page'

export default function Page() {
	return (
		<ProtectedRoute>
			<MainLayout>
				<StatsPage />
			</MainLayout>
		</ProtectedRoute>
	)
}
