import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { StudySessionPage } from '@/src/presentation/components/pages/study-session-page'

export default async function Page({ params }: { params: Promise<{ id: string }> }) {
	const { id } = await params
	return (
		<ProtectedRoute>
			<MainLayout>
				<StudySessionPage deckId={id} />
			</MainLayout>
		</ProtectedRoute>
	)
}
