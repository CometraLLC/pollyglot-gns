import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { TranslatePage } from '@/src/presentation/components/pages/translate-page'

export default function Page() {
	return (
		<ProtectedRoute>
			<MainLayout>
				<TranslatePage />
			</MainLayout>
		</ProtectedRoute>
	)
}
