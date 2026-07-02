import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { StudyPage } from '@/src/presentation/components/pages/study-page'

export default function Page() {
	return (
		<ProtectedRoute>
			<MainLayout>
				<StudyPage />
			</MainLayout>
		</ProtectedRoute>
	)
}
